package worker

import (
	"fmt"
	"io/ioutil"

	globalfunc "github.com/litekube/LiteKube/pkg/global"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
	"github.com/litekube/LiteKube/pkg/options/worker/global"
	"github.com/litekube/LiteKube/pkg/options/worker/kubelet"
	"github.com/litekube/LiteKube/pkg/options/worker/kubeproxy"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

type WorkerOptions struct {
	GlobalOptions     *global.GlobalOptions         `yaml:"global"`
	KubeletOptions    *kubelet.KubeletOptions       `yaml:"kubelet"`
	KubeProxyOptions  *kubeproxy.KubeProxyOptions   `yaml:"kube-proxy"`
	NetmamagerOptions *netmanager.NetManagerOptions `yaml:"network-manager"`
}

func NewWorkerOptions() *WorkerOptions {
	return &WorkerOptions{
		GlobalOptions:     global.NewGlobalOptions(),
		KubeletOptions:    kubelet.NewKubeletOptions(),
		KubeProxyOptions:  kubeproxy.NewKubeProxyOptions(),
		NetmamagerOptions: netmanager.NewNetManagerOptions(),
	}
}

func (opt *WorkerOptions) LoadConfig() error {
	// use default config
	if len(ConfigFile) < 1 {
		klog.Warningf("you did not specify a value for --%s=%s, the program will start with the default value. Use --help for more information", ConfigFileFlagName, ConfigFile)
		return nil
	}

	if !globalfunc.Exists(ConfigFile) {
		klog.Warningf("config file specify by --%s=%s is not exist, we will ignore this parameter. Use --help for more information", ConfigFileFlagName, ConfigFile)
		return nil
	}
	// try to read config file
	bytes, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return fmt.Errorf("unexpected errors while reading config from file specify by --%s=%s. Use --help for more information. Err: %v", ConfigFileFlagName, ConfigFile, err)
	}

	// unmarshal yaml
	if err := yaml.Unmarshal(globalfunc.ReplaceCurrent(globalfunc.ReplaceHome(bytes)), opt); err != nil {
		return fmt.Errorf("error while unmarshal config from file specify by --%s=%s, error info: %s", ConfigFileFlagName, ConfigFile, err.Error())
	}

	if err := opt.KubeletOptions.CheckReservedOptions(); err != nil {
		return err
	}

	if err := opt.KubeProxyOptions.CheckReservedOptions(); err != nil {
		return err
	}

	return nil
}

// add yaml format help tips
func (opt *WorkerOptions) ConfigHelpSection() []*help.Section {
	return []*help.Section{
		opt.GlobalOptions.HelpSection(),
		opt.KubeletOptions.HelpSection(),
		opt.KubeProxyOptions.HelpSection(),
		opt.NetmamagerOptions.HelpSection(),
	}
}

// add flags help tips
func (opt *WorkerOptions) HelpSections() []*help.Section {
	flagSection := help.NewSection("FLAGS", "", nil)
	flagSection.AddTip("--"+ConfigFileFlagName, "string", "YAML File to store leader startup parameters", "")
	flagSection.AddTip("--join", "string", "[only need for the first time]. string use to join to cluster, you can get from leader by seek administrator for help.", "")
	flagSection.AddTip("--versions", "string", "view the version info, value: {true,false,simple,raw}. ", "false")
	flagSection.AddTip("--help", "bool", "print usage", "false")

	return []*help.Section{flagSection}
}

func (opt *WorkerOptions) PrintFlags(printFunc func(format string, a ...interface{}) error) error {
	printFunc("[flags]:")
	opt.GlobalOptions.PrintFlags("global", printFunc)
	opt.KubeletOptions.PrintFlags("kubelet", printFunc)
	opt.KubeProxyOptions.PrintFlags("kube-proxy", printFunc)
	opt.NetmamagerOptions.PrintFlags("network-manager", printFunc)
	return nil
}
