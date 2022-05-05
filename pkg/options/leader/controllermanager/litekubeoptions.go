package controllermanager

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// options for Litekube to start kube-controller-manager
type ControllerManagerLitekubeOptions struct {
	AllocateNodeCidrs            bool   `yaml:"allocate-node-cidrs"`
	ClusterCidr                  string `yaml:"cluster-cidr"`
	Profiling                    bool   `yaml:"profiling"`
	UseServiceAccountCredentials bool   `yaml:"use-service-account-credentials"`
}

var DefaultCMLO ControllerManagerLitekubeOptions = ControllerManagerLitekubeOptions{
	AllocateNodeCidrs:            false,
	ClusterCidr:                  "172.17.0.0/16",
	Profiling:                    false,
	UseServiceAccountCredentials: true,
}

func NewControllerManagerLitekubeOptions() *ControllerManagerLitekubeOptions {
	options := DefaultCMLO
	return &options
}

func (opt *ControllerManagerLitekubeOptions) AddTips(section *help.Section) {
	section.AddTip("allocate-node-cidrs", "bool", "Should CIDRs for Pods be allocated and set on the cloud provider.", fmt.Sprintf("%t", DefaultCMLO.AllocateNodeCidrs))
	section.AddTip("cluster-cidr", "string", "CIDR Range for Pods in cluster. Requires --allocate-node-cidrs to be true", DefaultCMLO.ClusterCidr)
	section.AddTip("profiling", "bool", "Enable profiling via web interface host:port/debug/pprof/", fmt.Sprintf("%t", DefaultCMLO.Profiling))
	section.AddTip("use-service-account-credentials", "bool", "If true, use individual service account credentials for each controller.", fmt.Sprintf("%t", DefaultCMLO.UseServiceAccountCredentials))

}
