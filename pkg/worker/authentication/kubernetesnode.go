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
	"time"

	"github.com/litekube/LiteKube/pkg/global"
	"github.com/litekube/LiteKube/pkg/grpcclient"
	leaderruntime "github.com/litekube/LiteKube/pkg/leader/runtime"
	"github.com/litekube/LiteKube/pkg/leader/runtime/control"
	"github.com/litekube/LiteKube/pkg/likutemplate"
	globaloptions "github.com/litekube/LiteKube/pkg/options/worker/global"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

type KubernetesNode struct {
	RootCertDir             string
	CertDir                 string
	KubeletDir              string
	KubeproxyDir            string
	KubeProxyKubeConfig     string
	BootStrapKubeConfig     string
	KubeletConfig           string
	LeaderIp                string
	LeaderPort              uint16
	BootStrapToken          string
	LeaderNodeToken         string
	CurrentNodeToken        string
	RawToken                string
	ValidateApiserverCAFile string
	kubeletServerCert       string
	kubeletServerKey        string
	AdditionFile            string
	RegisterClient          *leaderruntime.NetWorkRegisterClient
	ClusterDNS              string
}

func NewKubernetesNode(rootCertPath string, leaderToken string, registerClient *leaderruntime.NetWorkRegisterClient) *KubernetesNode {
	tokens := strings.SplitN(leaderToken, "@", 2)
	if len(tokens[0]) < 1 || len(tokens[1]) < 1 {
		klog.Errorf("bad leader-token format")
		return nil
	}

	leaderNodeToken := tokens[0]
	bootstrapToken := tokens[1]

	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls")
	}
	certDir := filepath.Join(rootCertPath, leaderToken)

	leaderIp, err := registerClient.QueryIpByToken(leaderNodeToken)
	if err != nil || leaderIp == "" {
		klog.Errorf("fail to query ip of leader")
		return nil
	}

	kubeletDir := filepath.Join(certDir, "kubelet")
	kubeproxyDir := filepath.Join(certDir, "kube-proxy")
	if err := os.MkdirAll(kubeletDir, os.FileMode(0644)); err != nil {
		klog.Errorf("fail to create directory: %s", kubeletDir)
		return nil
	}
	if err := os.MkdirAll(kubeproxyDir, os.FileMode(0644)); err != nil {
		klog.Errorf("fail to create directory: %s", kubeproxyDir)
		return nil
	}

	return &KubernetesNode{
		RootCertDir:             rootCertPath,
		CertDir:                 certDir,
		KubeletDir:              kubeletDir,
		KubeproxyDir:            kubeproxyDir,
		KubeProxyKubeConfig:     filepath.Join(kubeproxyDir, "kube-proxy.kubeconfig"),
		BootStrapKubeConfig:     filepath.Join(kubeletDir, "bootstrap.kubeconfig"),
		KubeletConfig:           filepath.Join(kubeletDir, "kubelet.config"),
		ValidateApiserverCAFile: filepath.Join(kubeletDir, "validate-ca.crt"),
		kubeletServerCert:       filepath.Join(kubeletDir, "kubelet-server.crt"),
		kubeletServerKey:        filepath.Join(kubeletDir, "kubelet-server.key"),
		AdditionFile:            filepath.Join(certDir, "addition.map"),
		LeaderIp:                leaderIp,
		LeaderPort:              6442,
		BootStrapToken:          bootstrapToken,
		LeaderNodeToken:         leaderNodeToken,
		CurrentNodeToken:        registerClient.NodeToken,
		RegisterClient:          registerClient,
	}
}

func (kn *KubernetesNode) GenerateOrSkip() error {
	if kn == nil {
		return fmt.Errorf("nil kubernetes node")
	}

	if kn.BootStrapToken == "" || kn.LeaderNodeToken == "" {
		return fmt.Errorf("bad token to bootstrap for worker-kubernetes")
	}

	if kn.Check() {
		return nil
	}

	return kn.TLSBootStrap()
}

func (kn *KubernetesNode) Check() bool {
	return global.Exists(kn.KubeProxyKubeConfig, kn.BootStrapKubeConfig, kn.KubeletConfig, kn.ValidateApiserverCAFile, kn.AdditionFile)
}

func (kn *KubernetesNode) TLSBootStrap() error {
	if kn == nil {
		return fmt.Errorf("nil kubernetes node")
	}

	bootstrapToken := kn.BootStrapToken
	if bootstrapToken == "local" {
		bytes, err := ioutil.ReadFile(filepath.Join(global.HomePath, ".litekube/bootstrap-token"))
		if err != nil {
			return fmt.Errorf("leader-token=%s@local is only allow while running with leader in same node", global.ReservedNodeToken)
		}
		bootstrapToken = string(bytes)
	}
	auth := &grpcclient.TokenAuthentication{
		Token: bootstrapToken,
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
		if value.GetStatusCode() < 200 || value.GetStatusCode() > 299 {
			klog.Errorf("fail to bootstrap for kube-proxy, err: %s", value.GetMessage())
			return fmt.Errorf("fail to bootstrap for kube-proxy, err: %s", value.GetMessage())
		}

		kn.WriteAddition(map[string]string{"cluster-cidr": value.GetClusterCIDR()})
		if bytes, err := base64.StdEncoding.DecodeString(value.GetKubeconfig()); err != nil {
			return err
		} else {
			ioutil.WriteFile(kn.KubeProxyKubeConfig, bytes, os.FileMode(0644))
		}
	}

	// bootstrap for kubelet
	if value, err := client.BootStrapKubelet(context.Background(), &control.BootStrapKubeletRequest{NodeToken: kn.CurrentNodeToken}); err != nil {
		klog.Errorf("fail to bootstrap kubelet certificates")
		return err
	} else {
		if value.GetStatusCode() < 200 || value.GetStatusCode() > 299 {
			klog.Errorf("fail to bootstrap for kubelet, err: %s", value.GetMessage())
			return fmt.Errorf("fail to bootstrap for kubelet, err: %s", value.GetMessage())
		}

		// write bootstrap-kubeconfig
		if bytes, err := base64.StdEncoding.DecodeString(value.GetKubeconfig()); err != nil {
			return err
		} else {
			if err := ioutil.WriteFile(kn.BootStrapKubeConfig, bytes, os.FileMode(0644)); err != nil {
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

		// write server-cert
		if bytes, err := base64.StdEncoding.DecodeString(value.GetServerCert()); err != nil {
			return err
		} else {
			if err := ioutil.WriteFile(kn.kubeletServerCert, bytes, os.FileMode(0644)); err != nil {
				return err
			}
		}

		// write server-key
		if bytes, err := base64.StdEncoding.DecodeString(value.GetServerKey()); err != nil {
			return err
		} else {
			if err := ioutil.WriteFile(kn.kubeletServerKey, bytes, os.FileMode(0644)); err != nil {
				return err
			}
		}

		// write validate-ca
		kn.ClusterDNS = value.GetClusterDNS()

		nodeIp, err := kn.RegisterClient.QueryIp()
		if err != nil {
			nodeIp = "0.0.0.0"
		}

		// write default-kubelet config
		buf := &bytes.Buffer{}
		data := struct {
			CaPath          string
			KubeletServerIp string
			ServerCertPath  string
			ServerKeyPath   string
			CluserDNS       string
		}{
			CaPath:          kn.ValidateApiserverCAFile,
			KubeletServerIp: nodeIp,
			ServerCertPath:  kn.kubeletServerCert,
			ServerKeyPath:   kn.kubeletServerKey,
			CluserDNS:       kn.ClusterDNS,
		}

		likutemplate.Kubelet_config_template.Execute(buf, &data)
		if err := ioutil.WriteFile(kn.KubeletConfig, buf.Bytes(), os.FileMode(0644)); err != nil {
			return err
		}
		kn.WriteAddition(map[string]string{"cluster-dns": kn.ClusterDNS})
	}

	if !kn.Check() {
		return fmt.Errorf("bootstrap for worker meet unknow error")
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

func (kn *KubernetesNode) ReadAllAddition() (map[string]string, error) {
	if !global.Exists(kn.AdditionFile) {
		return make(map[string]string), nil
	}

	bytes, err := ioutil.ReadFile(kn.AdditionFile)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (kn *KubernetesNode) ClearAddition() error {
	if global.Exists(kn.AdditionFile) {
		return os.Remove(kn.AdditionFile)
	}

	return nil
}

func (kn *KubernetesNode) WriteAddition(new_data map[string]string) error {
	if new_data == nil {
		return nil
	}

	if err := os.MkdirAll(kn.CertDir, os.FileMode(0644)); err != nil {
		return err
	}

	data, err := kn.ReadAllAddition()
	if err != nil {
		return err
	}

	for key, value := range new_data {
		data[key] = value
	}

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(kn.AdditionFile, bytes, fs.FileMode(0644))
}
