package scheduler

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// Empirically assigned parameters are not recommended
type SchedulerProfessionalOptions struct {
	BindAddress              string `yaml:"bind-address"`
	SecurePort               uint16 `yaml:"secure-port"`
	LeaderElect              bool   `yaml:"leader-elect"`
	KubeConfig               string `yaml:"kubeconfig"`
	AuthorizationKubeconfig  string `yaml:"authorization-kubeconfig"`
	AuthenticationKubeconfig string `yaml:"authentication-kubeconfig"`
}

var DefaultSPO SchedulerProfessionalOptions = SchedulerProfessionalOptions{
	BindAddress: "0.0.0.0",
	SecurePort:  10259,
	LeaderElect: false,
}

func NewSchedulerProfessionalOptions() *SchedulerProfessionalOptions {
	options := DefaultSPO
	return &options
}

func (opt *SchedulerProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port. ", DefaultSPO.BindAddress)
	section.AddTip("secure-port", "uint16", "The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all.", fmt.Sprintf("%d", DefaultSPO.SecurePort))
	section.AddTip("leader-elect", "bool", "Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.", fmt.Sprintf("%t", DefaultSPO.LeaderElect))
	section.AddTip("authorization-kubeconfig", "string", "kubeconfig file pointing at the 'core' kubernetes server with enough rights to create subjectaccessreviews.authorization.k8s.io. ", DefaultSPO.AuthorizationKubeconfig)
	section.AddTip("authentication-kubeconfig", "string", "kubeconfig file pointing at the 'core' kubernetes server with enough rights to create tokenreviews.authentication.k8s.io.", DefaultSPO.AuthenticationKubeconfig)
	section.AddTip("kubeconfig", "string", "deprecated. Path to kubeconfig file with authorization and master location information. ", DefaultSPO.KubeConfig)
}
