package leader

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/leader/apiserver"
	"github.com/litekube/LiteKube/pkg/options/leader/global"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
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

func (opt *LeaderOptions) LoadConfig() error {
	// use default config
	if len(ConfigFile) < 1 {
		klog.Warningf("you did not specify a value for --%s=%s, the program will start with the default value. Use --help for more information", ConfigFileFlagName, ConfigFile)
		return nil
	}

	// try to read config file
	bytes, err := ioutil.ReadFile(ConfigFile)
	if err == os.ErrNotExist {
		klog.Warningf("config file specify by --%s=%s is not exist, we will ignore this parameter. Use --help for more information", ConfigFileFlagName, ConfigFile)
		return nil
	} else if err != nil {
		return fmt.Errorf("unexpected errors while reading config from file specify by --%s=%s. Use --help for more information", ConfigFileFlagName, ConfigFile)
	}

	// unmarshal yaml
	if err := yaml.Unmarshal(bytes, opt); err != nil {
		return fmt.Errorf("error while unmarshal config from file specify by --%s=%s, error info: %s", ConfigFileFlagName, ConfigFile, err.Error())
	}

	return nil
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
	flagSection.AddTip("--help", "bool", "print usage", "false")

	return []*help.Section{flagSection}
}

func (opt *LeaderOptions) PrintFlags(printFunc func(format string, a ...interface{}) error) error {
	printFunc("[flags]:")
	opt.GlobalOptions.PrintFlags("litekube", printFunc)
	opt.ApiserverOptions.PrintFlags("kube-apiserver", printFunc)
	return nil
}
