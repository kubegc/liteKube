package global

import (
	"fmt"
	"path/filepath"

	"github.com/litekube/LiteKube/pkg/global"
	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type GlobalOptions struct {
	WorkDir     string `yaml:"work-dir"`
	LogDir      string `yaml:"log-dir"`
	LogToDir    bool   `yaml:"log-to-dir"`
	LogToStd    bool   `yaml:"log-to-std"`
	LeaderToken string `yaml:"leader-token"`
}

var DefaultGO GlobalOptions = GlobalOptions{
	WorkDir:     filepath.Join(global.HomePath, ".litekube/"),
	LogDir:      "",
	LogToStd:    true,
	LogToDir:    true,
	LeaderToken: global.ReservedNodeToken,
}

func NewGlobalOptions() *GlobalOptions {
	options := DefaultGO
	return &options
}

func (opt *GlobalOptions) HelpSection() *help.Section {
	section := help.NewSection("global", "leader startup parameters and common args for kubernetes components", nil)
	section.AddTip("work-dir", "string", "dir to store file generate by litekube", DefaultGO.WorkDir)
	section.AddTip("log-dir", "string", "fold path to store logs", "$WorkDir/logs")
	section.AddTip("log-to-dir", "bool", "store log to disk or not", fmt.Sprintf("%t", DefaultGO.LogToDir))
	section.AddTip("log-to-std", "bool", "print log to disk or not", fmt.Sprintf("%t", DefaultGO.LogToStd))
	section.AddTip("leader-token", "string", "token to join into k8s cluster", DefaultGO.LeaderToken)
	//section.AddTip("worker-config", "string", "worker config, --enable-work=true is recommanded", DefaultGO.WorkerConfig)
	return section
}

// print all flags
func (opt *GlobalOptions) PrintFlags(prefix string, printFunc func(format string, a ...interface{}) error) error {
	// print flags
	flags, err := common.StructToMap(opt)
	if err != nil {
		return err
	}
	common.PrintMap(flags, prefix, printFunc)
	return nil
}
