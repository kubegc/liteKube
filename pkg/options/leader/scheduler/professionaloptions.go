package scheduler

import "github.com/litekube/LiteKube/pkg/help"

// Empirically assigned parameters are not recommended
type SchedulerProfessionalOptions struct {
	BindAddress string `yaml:"bind-address"`
	SecurePort  int16  `yaml:"secure-port"`
	LeaderElect bool   `yaml:"leader-elect"`
}

func NewSchedulerProfessionalOptions() *SchedulerProfessionalOptions {
	return &SchedulerProfessionalOptions{}
}

func (opt *SchedulerProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port. ", "0.0.0.0")
	section.AddTip("secure-port", "int16", "The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all.", "10259")
	section.AddTip("leader-elect", "bool", "Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.", "false")
}
