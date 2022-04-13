package global

import (
	"fmt"
	"os/user"
	"path/filepath"
	"sort"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type PrintFunc func(format string, a ...interface{}) error

type GlobalOptions struct {
	WorkDir      string `yaml:"work-dir"`
	LogDir       string `yaml:"log-dir"`
	LogToDir     bool   `yaml:"log-to-dir"`
	LogToStd     bool   `yaml:"log-to-std"`
	EnableWorker bool   `yaml:"enable-worker"`
	WorkerConfig string `yaml:"worker-config"`
}

var defaultGO GlobalOptions = GlobalOptions{
	WorkDir:      filepath.Join(GetHomeDir(), "litekube/"),
	LogDir:       "/var/log/litekube/",
	LogToStd:     true,
	LogToDir:     false,
	EnableWorker: false,
}

func GetHomeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		return "/"
	}

	return currentUser.HomeDir
}

func NewGlobalOptions() *GlobalOptions {
	options := defaultGO
	return &options
}

func (opt *GlobalOptions) HelpSection() *help.Section {
	section := help.NewSection("global", "leader startup parameters and common args for kubernetes components", nil)

	section.AddTip("work-dir", "string", "dir to store file generate by litekube", defaultGO.WorkDir)
	section.AddTip("log-dir", "string", "fold path to store logs", defaultGO.LogDir)
	section.AddTip("log-to-dir", "bool", "store log to disk or not", fmt.Sprintf("%t", defaultGO.LogToDir))
	section.AddTip("log-to-std", "bool", "print log to disk or not", fmt.Sprintf("%t", defaultGO.LogToStd))
	section.AddTip("enable-worker", "bool", "run worker together or not", fmt.Sprintf("%t", defaultGO.EnableWorker))
	section.AddTip("worker-config", "string", "worker config, --enable-work=true is recommanded", defaultGO.WorkerConfig)
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
