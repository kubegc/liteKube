package global

import (
	"sort"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type PrintFunc func(format string, a ...interface{}) error

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

// print all flags
func (opt *GlobalOptions) PrintFlags(prefix string, printFunc func(format string, a ...interface{}) error) error {
	// print flags
	flags, err := common.StructToMap(opt)
	if err != nil {
		return err
	}
	printMap(flags, prefix, printFunc)
	return nil
}

func printMap(m map[string]string, prefix string, printFunc PrintFunc) {
	if m == nil {
		return
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		printFunc("--%s-%s=%s", prefix, key, m[key])
	}
}
