package global

import "github.com/litekube/LiteKube/pkg/help"

type GlobalOptions struct {
	LogDir   string `yaml:"log-dir"`
	LogToDir bool   `yaml:"log-to-dir"`
	LogToStd bool   `yaml:"log-to-std"`
}

func NewGlobalOptions() *GlobalOptions {
	return &GlobalOptions{}
}

func (opt *GlobalOptions) HelpSection() *help.Section {
	section := help.NewSection("global", "leader startup parameters and common args for kubernetes components", nil)

	section.AddTip("log-dir", "string", "fold path to store logs", "")
	section.AddTip("log-to-dir", "bool", "store log to disk or not", "false")
	section.AddTip("log-to-std", "bool", "print log to disk or not", "true")
	return section
}
