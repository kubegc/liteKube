package leader

import (
	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/leader/apiserver"
)

type LeaderOptions struct {
	ApiserverOptions *apiserver.ApiserverOptions `yaml:"kube-apiserver"`
}

func NewLeaderOptions() *LeaderOptions {
	return &LeaderOptions{
		ApiserverOptions: apiserver.NewApiserverOptions(),
	}
}

func (opt *LeaderOptions) ConfigHelpSection() *help.Section {
	// add tips for apiserver
	leaderSection := help.NewSection("Leader", "setting for kube-apiserver,kube-controller-manager,kube-scheduler and litekube additions", nil)
	leaderSection.AddSection(opt.ApiserverOptions.HelpSection())

	return leaderSection
}

func (opt *LeaderOptions) HelpSections() []*help.Section {
	// add tips for apiserver
	flagSection := help.NewSection("FLAGS", "", nil)
	flagSection.AddTip("--config-file", "string", "YAML File to store leader startup parameters", "")
	flagSection.AddTip("--version", "string", "view the version info, value: {true,false,simple,raw}. ", "false")

	return []*help.Section{flagSection}
}
