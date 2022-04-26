package authentication

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/litekube/LiteKube/pkg/global"
	leaderruntime "github.com/litekube/LiteKube/pkg/leader/runtime"
	globaloptions "github.com/litekube/LiteKube/pkg/options/worker/global"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

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
	if err != nil {
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

func (kn *KubernetesNode) GenerateOrSkip(address string, port int) error {
	if kn == nil {
		return fmt.Errorf("nil kubernetes node")
	}

	return kn.TLSBootStrap()
}

func (kn *KubernetesNode) TLSBootStrap() error {
	if kn == nil {
		return fmt.Errorf("nil kubernetes node")
	}

	if global.Exists(kn.KubeProxyKubeConfig, kn.BootStrapKubeConfig, kn.KubeProxyKubeConfig, kn.ValidateApiserverCAFile, kn.AdditionFile) {
		return nil
	} else {
		// try to download certificates here
		return nil
	}
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
