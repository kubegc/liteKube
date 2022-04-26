package kubeproxy

import (
	"github.com/litekube/LiteKube/pkg/help"
)

// Empirically assigned parameters are not recommended
type KubeProxyProfessionalOptions struct {
	HostnameOverride string `yaml:"hostname-override"`
	ClusterCidr      string `yaml:"cluster-cidr"`
	ProxyMode        string `yaml:"proxy-mode"`
	Kubeconfig       string `yaml:"kubeconfig"`
}

func NewKubeProxyProfessionalOptions() *KubeProxyProfessionalOptions {
	options := DefaultKPPO
	return &options
}

var DefaultKPPO KubeProxyProfessionalOptions = KubeProxyProfessionalOptions{
	ProxyMode: "ipvs",
}

func (opt *KubeProxyProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("hostname-override", "string", "If non-empty, will use this string as identification instead of the actual hostname.", DefaultKPPO.HostnameOverride)
	section.AddTip("cluster-cidr", "string", "The CIDR range of pods in the cluster.", DefaultKPPO.ClusterCidr)
	section.AddTip("proxy-mode", "string", "Which proxy mode to use: 'userspace' (older) or 'iptables' (faster) or 'ipvs' or 'kernelspace' (windows).", DefaultKPPO.ProxyMode)
	section.AddTip("kubeconfig", "string", "Path to kubeconfig file with authorization information.", DefaultKPPO.Kubeconfig)
}
