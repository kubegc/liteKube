package kubeproxy

import (
	"github.com/litekube/LiteKube/pkg/help"
)

// options for Litekube to start kube-controller-manager
type KubeProxyLitekubeOptions struct {
}

var DefaultKPLO KubeProxyLitekubeOptions = KubeProxyLitekubeOptions{}

func NewKubeProxyLitekubeOptions() *KubeProxyLitekubeOptions {
	options := DefaultKPLO
	return &options
}

func (opt *KubeProxyLitekubeOptions) AddTips(section *help.Section) {
}
