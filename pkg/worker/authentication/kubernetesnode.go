package authentication

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/litekube/LiteKube/pkg/global"
	"github.com/litekube/LiteKube/pkg/grpcclient"
	leaderruntime "github.com/litekube/LiteKube/pkg/leader/runtime"
	"github.com/litekube/LiteKube/pkg/leader/runtime/control"
	globaloptions "github.com/litekube/LiteKube/pkg/options/worker/global"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var kubelet_config_template = template.Must(template.New("kubelet_config").Parse(`kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  x509:
    clientCAFile: "{{.CaPath}}"
  webhook:
    enabled: true
    cacheTTL: 2m0s
  anonymous:
    enabled: false
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 5m0s
    cacheUnauthorizedTTL: 30s
address: "0.0.0.0"
port: 10250
readOnlyPort: 10255
cgroupDriver: systemd
hairpinMode: promiscuous-bridge
serializeImagePulls: false
clusterDomain: cluster.local.
clusterDNS:
- "{{.CluserDNS}}"
`))

type KubernetesNode struct {
	CertDir                 string
	KubeProxyKubeConfig     string
	BootStrapKubeConfig     string
	KubeletConfig           string
	LeaderIp                string
	LeaderPort              uint16
	BootStrapToken          string
	LeaderNodeToken         string
	RawToken                string
	ValidateApiserverCAFile string
	AdditionFile            string
	RegisterClient          *leaderruntime.NetWorkRegisterClient
	CluserDNS               string
}

func NewKubernetesNode(rootCertPath string, leaderToken string, registerClient *leaderruntime.NetWorkRegisterClient) *KubernetesNode {
	tokens := strings.SplitN(leaderToken, "@", 2)
	if len(tokens[0]) < 1 || len(tokens[1]) < 1 {
		klog.Errorf("bad leader-token format")
		return nil
	}

	leaderNodeToken := tokens[0]
	bootstrapToken := tokens[1]

	certDir := filepath.Join(rootCertPath, leaderToken)
	if rootCertPath == "" {
		certDir = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls", leaderToken)
	}

	leaderIp, err := registerClient.QueryIpByToken(leaderNodeToken)
	if err != nil || leaderIp == "" {
		klog.Errorf("fail to query ip of leader")
		return nil
	}

	return &KubernetesNode{
		CertDir:                 certDir,
		KubeProxyKubeConfig:     filepath.Join(certDir, "kube-proxy", "kube-proxy.kubeconfig"),
		BootStrapKubeConfig:     filepath.Join(certDir, "kubelet", "bootstrap.kubeconfig"),
		KubeletConfig:           filepath.Join(certDir, "kubelet", "kubelet.config"),
		ValidateApiserverCAFile: filepath.Join(certDir, "kubelet", "validate-ca.cert"),
		AdditionFile:            filepath.Join(certDir, "addition.map"),
		LeaderIp:                leaderIp,
		LeaderPort:              6440,
		BootStrapToken:          bootstrapToken,
		LeaderNodeToken:         leaderNodeToken,
		RegisterClient:          registerClient,
	}
}

func (kn *KubernetesNode) GenerateOrSkip() error {
	if kn == nil {
		return fmt.Errorf("nil kubernetes node")
	}

	return kn.TLSBootStrap()
}

func (kn *KubernetesNode) TLSBootStrap() error {
	if kn == nil {
		return fmt.Errorf("nil kubernetes node")
	}

	if global.Exists(kn.KubeProxyKubeConfig, kn.BootStrapKubeConfig, kn.KubeletConfig, kn.ValidateApiserverCAFile, kn.AdditionFile) {
		return nil
	} else {
		auth := &grpcclient.TokenAuthentication{
			Token: kn.BootStrapToken,
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", kn.LeaderIp, kn.LeaderPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(auth))
		if err != nil {
			return err
		}

		client := control.NewLeaderControlClient(conn)

		for i := 0; i < 100; i++ {
			if _, err := client.CheckHealth(context.Background(), &control.NoneValue{}); err != nil {
				klog.Warningf("waiting for leader controller to start, try %d/100 times", i)
				time.Sleep(1 * time.Second)
				continue
			} else {
				break
			}
		}

		// bootstrap for kube-proxy
		if value, err := client.BootStrapKubeProxy(context.Background(), &control.BootStrapKubeProxyRequest{}); err != nil {
			klog.Errorf("fail to bootstrap kube-proxy certificates")
			return err
		} else {
			kn.WriteAddition(map[string]string{"cluster-cidr": value.GetClusterCIDR()})
			if bytes, err := base64.StdEncoding.DecodeString(value.GetKubeconfig()); err != nil {
				return err
			} else {
				ioutil.WriteFile(kn.KubeProxyKubeConfig, bytes, os.FileMode(0644))
			}
		}

		// bootstrap for kubelet
		if value, err := client.BootStrapKubelet(context.Background(), &control.BootStrapKubeletRequest{}); err != nil {
			klog.Errorf("fail to bootstrap kubelet certificates")
			return err
		} else {
			// write bootstrap-kubeconfig
			if bytes, err := base64.StdEncoding.DecodeString(value.GetKubeconfig()); err != nil {
				return err
			} else {
				if err := ioutil.WriteFile(kn.KubeletConfig, bytes, os.FileMode(0644)); err != nil {
					return err
				}
			}

			// write validate-ca
			if bytes, err := base64.StdEncoding.DecodeString(value.GetValidataCaCert()); err != nil {
				return err
			} else {
				if err := ioutil.WriteFile(kn.ValidateApiserverCAFile, bytes, os.FileMode(0644)); err != nil {
					return err
				}
			}

			// write validate-ca
			kn.CluserDNS = value.GetClusterDNS()

			// write default-kubelet config
			buf := &bytes.Buffer{}
			data := struct {
				CaPath    string
				CluserDNS string
			}{
				CaPath:    kn.ValidateApiserverCAFile,
				CluserDNS: kn.CluserDNS,
			}

			kubelet_config_template.Execute(buf, &data)
			if err := ioutil.WriteFile(kn.KubeletConfig, buf.Bytes(), os.FileMode(0644)); err != nil {
				return err
			}
			kn.WriteAddition(map[string]string{"cluster-dns": kn.CluserDNS})
		}

	}

	return nil
}

func (kn *KubernetesNode) ReadAddition(key string) (string, error) {
	if !global.Exists(kn.AdditionFile) {
		return "", fmt.Errorf("fail to find addition file")
	}

	bytes, err := ioutil.ReadFile(kn.AdditionFile)
	if err != nil {
		return "", err
	}

	data := make(map[string]string)
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return "", err
	}

	if value, ok := data[key]; ok {
		return value, nil
	} else {
		return "", fmt.Errorf("bad key to read addition")
	}
}

func (kn *KubernetesNode) WriteAddition(data map[string]string) error {
	if data == nil {
		return nil
	}

	if err := os.MkdirAll(kn.CertDir, os.FileMode(0644)); err != nil {
		return err
	}

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(kn.AdditionFile, bytes, fs.FileMode(0644))
}
