package kubelet

import (
	"github.com/litekube/LiteKube/pkg/help"
)

// options for Litekube to start kube-controller-manager
type KubeletLitekubeOptions struct {
	PodInfraContainerImage string `yaml:"pod-infra-container-image"`
	CertDir                string `yaml:"cert-dir"`
}

var DefaultKLO KubeletLitekubeOptions = KubeletLitekubeOptions{
	PodInfraContainerImage: "registry.cn-hangzhou.aliyuncs.com/google-containers/pause-amd64:3.0",
}

func NewKubeletLitekubeOptions() *KubeletLitekubeOptions {
	options := DefaultKLO
	return &options
}

func (opt *KubeletLitekubeOptions) AddTips(section *help.Section) {
	section.AddTip("pod-infra-container-image", "string", "Specified image will not be pruned by the image garbage collector.", DefaultKLO.PodInfraContainerImage)
	section.AddTip("cert-dir", "string", "The directory where the TLS certs are located. ", DefaultKLO.CertDir)
}
