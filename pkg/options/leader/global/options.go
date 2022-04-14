package global

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/litekube/LiteKube/pkg/global"
	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type PrintFunc func(format string, a ...interface{}) error

type GlobalOptions struct {
	WorkDir      string `yaml:"work-dir"`
	LogDir       string `yaml:"log-dir"`
	LogToDir     bool   `yaml:"log-to-dir"`
	LogToStd     bool   `yaml:"log-to-std"`
	RunKine      bool   `yaml:"run-kine"`
	EnableWorker bool   `yaml:"enable-worker"`
	WorkerConfig string `yaml:"worker-config"`
}

var DefaultGO GlobalOptions = GlobalOptions{
	WorkDir:      filepath.Join(global.HomePath, "litekube/"),
	LogDir:       filepath.Join(global.HomePath, "litekube/logs/"),
	LogToStd:     true,
	LogToDir:     false,
	RunKine:      true,
	EnableWorker: false,
}

func NewGlobalOptions() *GlobalOptions {
	options := DefaultGO
	return &options
}

func (opt *GlobalOptions) HelpSection() *help.Section {
	section := help.NewSection("global", "leader startup parameters and common args for kubernetes components", nil)
	section.AddTip("work-dir", "string", "dir to store file generate by litekube", DefaultGO.WorkDir)
	section.AddTip("log-dir", "string", "fold path to store logs", DefaultGO.LogDir)
	section.AddTip("log-to-dir", "bool", "store log to disk or not", fmt.Sprintf("%t", DefaultGO.LogToDir))
	section.AddTip("log-to-std", "bool", "print log to disk or not", fmt.Sprintf("%t", DefaultGO.LogToStd))
	section.AddTip("run-kine", "bool", "run kine in leader process or not", fmt.Sprintf("%t", DefaultGO.RunKine))
	section.AddTip("enable-worker", "bool", "run worker together or not", fmt.Sprintf("%t", DefaultGO.EnableWorker))
	section.AddTip("worker-config", "string", "worker config, --enable-work=true is recommanded", DefaultGO.WorkerConfig)
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
