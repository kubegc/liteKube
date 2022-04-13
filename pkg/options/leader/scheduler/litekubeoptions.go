package scheduler

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// options for Litekube to start kube-scheduler
type SchedulerLitekubeOptions struct {
	Profiling                bool   `yaml:"profiling"`
	KubeConfig               string `yaml:"kubeconfig"`
	AuthorizationKubeconfig  string `yaml:"authorization-kubeconfig"`
	AuthenticationKubeconfig string `yaml:"authentication-kubeconfig"`
}

var defaultSLO SchedulerLitekubeOptions = SchedulerLitekubeOptions{
	Profiling: false,
}

func NewSchedulerLitekubeOptions() *SchedulerLitekubeOptions {
	options := defaultSLO
	return &options
}

func (opt *SchedulerLitekubeOptions) AddTips(section *help.Section) {
	section.AddTip("authorization-kubeconfig", "string", "kubeconfig file pointing at the 'core' kubernetes server with enough rights to create subjectaccessreviews.authorization.k8s.io. ", defaultSLO.AuthorizationKubeconfig)
	section.AddTip("authentication-kubeconfig", "string", "kubeconfig file pointing at the 'core' kubernetes server with enough rights to create tokenreviews.authentication.k8s.io.", defaultSLO.AuthenticationKubeconfig)
	section.AddTip("kubeconfig", "string", "deprecated. Path to kubeconfig file with authorization and master location information. ", defaultSLO.KubeConfig)
	section.AddTip("profiling", "bool", "deprecated. Enable profiling via web interface host:port/debug/pprof/.", fmt.Sprintf("%t", defaultSLO.Profiling))
}
