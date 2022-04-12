package leader

import (
	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/leader/apiserver"
	"github.com/litekube/LiteKube/pkg/options/leader/global"
)

type LeaderOptions struct {
	GlobalOptions    *global.GlobalOptions       `yaml:"global"`
	ApiserverOptions *apiserver.ApiserverOptions `yaml:"kube-apiserver"`
}

func NewLeaderOptions() *LeaderOptions {
	return &LeaderOptions{
		ApiserverOptions: apiserver.NewApiserverOptions(),
		GlobalOptions:    global.NewGlobalOptions(),
	}
}

// add yaml format help tips
func (opt *LeaderOptions) ConfigHelpSection() *help.Section {
	yamlSection := help.NewSection("Leader", "setting for kube-apiserver,kube-controller-manager,kube-scheduler and litekube additions", nil)
	// add for global
	yamlSection.AddSection(opt.GlobalOptions.HelpSection())

	// add for kube-apiserver
	yamlSection.AddSection(opt.ApiserverOptions.HelpSection())

	return yamlSection
}

// add flags help tips
func (opt *LeaderOptions) HelpSections() []*help.Section {
	// add tips for apiserver
	flagSection := help.NewSection("FLAGS", "", nil)
	flagSection.AddTip("--"+ConfigFileFlagName, "string", "YAML File to store leader startup parameters", "")
	flagSection.AddTip("--version", "string", "view the version info, value: {true,false,simple,raw}. ", "false")

	return []*help.Section{flagSection}
}
