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

var defaultCMLO ControllerManagerLitekubeOptions = ControllerManagerLitekubeOptions{
	AllocateNodeCidrs: false,
	Profiling:         false,
}

func NewControllerManagerLitekubeOptions() *ControllerManagerLitekubeOptions {
	options := defaultCMLO
	return &options
}

func (opt *ControllerManagerLitekubeOptions) AddTips(section *help.Section) {
	section.AddTip("allocate-node-cidrs", "bool", "Should CIDRs for Pods be allocated and set on the cloud provider.", fmt.Sprintf("%t", defaultCMLO.AllocateNodeCidrs))
	section.AddTip("cluster-cidr", "string", "CIDR Range for Pods in cluster. Requires --allocate-node-cidrs to be true", defaultCMLO.ClusterCidr)
	section.AddTip("profiling", "bool", "Enable profiling via web interface host:port/debug/pprof/", fmt.Sprintf("%t", defaultCMLO.Profiling))
	section.AddTip("use-service-account-credentials", "bool", "If true, use individual service account credentials for each controller.", fmt.Sprintf("%t", defaultCMLO.UseServiceAccountCredentials))

}
