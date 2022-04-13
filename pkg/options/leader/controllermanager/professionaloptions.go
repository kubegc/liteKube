package controllermanager

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// Empirically assigned parameters are not recommended
type ControllerManagerProfessionalOptions struct {
	BindAddress          string `yaml:"bind-address"`
	SecurePort           int16  `yaml:"secure-port"`
	LeaderElect          bool   `yaml:"leader-elect"`
	ConfigureCloudRoutes bool   `yaml:"configure-cloud-routes"`
	Controllers          string `yaml:"controllers"`
}

func NewControllerManagerProfessionalOptions() *ControllerManagerProfessionalOptions {
	options := defaultCMPO
	return &options
}

var defaultCMPO ControllerManagerProfessionalOptions = ControllerManagerProfessionalOptions{
	BindAddress:          "0.0.0.0",
	SecurePort:           10257,
	LeaderElect:          false,
	ConfigureCloudRoutes: false,
	Controllers:          "*,-service,-route,-cloud-node-lifecycle",
}

func (opt *ControllerManagerProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port. ", defaultCMPO.BindAddress)
	section.AddTip("secure-port", "int16", "The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all.", fmt.Sprintf("%d", defaultCMPO.SecurePort))
	section.AddTip("leader-elect", "bool", "Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.", fmt.Sprintf("%t", defaultCMPO.LeaderElect))
	section.AddTip("configure-cloud-routes", "bool", "Should CIDRs allocated by allocate-node-cidrs be configured on the cloud provider.", fmt.Sprintf("%t", defaultCMPO.ConfigureCloudRoutes))
	section.AddTip("controllers", "string", "A list of controllers to enable. ", defaultCMPO.Controllers)
}
