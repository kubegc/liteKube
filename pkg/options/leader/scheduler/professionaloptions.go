package scheduler

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// Empirically assigned parameters are not recommended
type SchedulerProfessionalOptions struct {
	BindAddress string `yaml:"bind-address"`
	SecurePort  int16  `yaml:"secure-port"`
	LeaderElect bool   `yaml:"leader-elect"`
}

var defaultSPO SchedulerProfessionalOptions = SchedulerProfessionalOptions{
	BindAddress: "0.0.0.0",
	SecurePort:  10259,
	LeaderElect: false,
}

func NewSchedulerProfessionalOptions() *SchedulerProfessionalOptions {
	options := defaultSPO
	return &options
}

func (opt *SchedulerProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port. ", defaultSPO.BindAddress)
	section.AddTip("secure-port", "int16", "The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all.", fmt.Sprintf("%d", defaultSPO.SecurePort))
	section.AddTip("leader-elect", "bool", "Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.", fmt.Sprintf("%t", defaultSPO.LeaderElect))
}
