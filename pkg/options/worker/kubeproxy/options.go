package kubeproxy

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type PrintFunc func(format string, a ...interface{}) error

// struct to store args from input
type KubeProxyOptions struct {
	ReservedOptions     map[string]string             `yaml:"reserve"`
	ProfessionalOptions *KubeProxyProfessionalOptions `yaml:"professional"`
	Options             *KubeProxyLitekubeOptions     `yaml:"options"`
	IgnoreOptions       map[string]string             `yaml:"-"`
}

func NewKubeProxyOptions() *KubeProxyOptions {
	return &KubeProxyOptions{
		ReservedOptions:     make(map[string]string),
		ProfessionalOptions: NewKubeProxyProfessionalOptions(),
		Options:             NewKubeProxyLitekubeOptions(),
		IgnoreOptions:       make(map[string]string),
	}
}

// delete keys already be disable or define in other block
func (opt *KubeProxyOptions) CheckReservedOptions() error {
	// deep copy options
	optionsMap, oErr := common.StructToMap(opt.Options)
	if oErr != nil {
		return oErr
	}

	for k := range optionsMap {
		if value, ok := opt.ReservedOptions[k]; ok {
			opt.IgnoreOptions[k] = value
			delete(opt.ReservedOptions, k)
		}
	}

	// deep copy professional-options
	professionalOptionsMap, pErr := common.StructToMap(opt.ProfessionalOptions)
	if pErr != nil {
		return pErr
	}

	for k := range professionalOptionsMap {
		if value, ok := opt.ReservedOptions[k]; ok {
			opt.IgnoreOptions[k] = value
			delete(opt.ReservedOptions, k)
		}
	}
	return nil
}

func (opt *KubeProxyOptions) HelpSection() *help.Section {
	section := help.NewSection("kube-proxy", "kube-proxy's startup parameters, we keep almost the same Settings as the original except logs relation: https://kubernetes.io/docs/reference/command-line-tools-reference/kube-proxy/", nil)

	reserveSection := help.NewSection("reserve", "reserve parameters, you can still add args unmentioned before refer to kube-proxy official website.", nil)
	reserveSection.AddTip("<name-1>", "<value-1>", "it will be parse as \"--<name-1>=<value-1>\"", "")
	reserveSection.AddTip("<name-n>", "<value-n>", "and so on", "")
	section.AddSection(reserveSection)

	professionalSection := help.NewSection("professional", "parameters are not recommended to set by users", nil)
	opt.ProfessionalOptions.AddTips(professionalSection)
	section.AddSection(professionalSection)

	litekubeoptionsSection := help.NewSection("options", "Litekube normal options", nil)
	opt.Options.AddTips(litekubeoptionsSection)
	section.AddSection(litekubeoptionsSection)

	return section
}

func (opt *KubeProxyOptions) ToMap() (map[string]string, error) {
	// check error define for flags
	opt.CheckReservedOptions()

	args := make(map[string]string)

	// deep copy reserved-options
	for k, v := range opt.ReservedOptions {
		args[k] = v
	}

	// deep copy options
	optionsMap, oErr := common.StructToMap(opt.Options)
	if oErr != nil {
		return nil, oErr
	}

	for k, v := range optionsMap {
		args[k] = v
	}

	// deep copy professional-options
	professionalOptionsMap, pErr := common.StructToMap(opt.ProfessionalOptions)
	if pErr != nil {
		return nil, pErr
	}

	for k, v := range professionalOptionsMap {
		args[k] = v
	}

	return args, nil
}

// print all flags
func (opt *KubeProxyOptions) PrintFlags(prefix string, printFunc func(format string, a ...interface{}) error) error {
	// print flags
	flags, err := opt.ToMap()
	if err != nil {
		return err
	}
	common.PrintMap(flags, prefix, printFunc)
	// print flags-ignored
	common.PrintMap(opt.IgnoreOptions, fmt.Sprintf("%s-<%s>", prefix, UnreserveTip), printFunc)
	return nil
}
