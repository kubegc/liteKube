package scheduler

import (
	"fmt"
	"sort"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/common"
)

type PrintFunc func(format string, a ...interface{}) error

// struct to store args from input
type SchedulerOptions struct {
	ReservedOptions     map[string]string             `yaml:"reserve"`
	ProfessionalOptions *SchedulerProfessionalOptions `yaml:"professional"`
	Options             *SchedulerLitekubeOptions     `yaml:"options"`
	IgnoreOptions       map[string]string             `yaml:"-"`
}

func NewSchedulerOptions() *SchedulerOptions {
	return &SchedulerOptions{
		ReservedOptions:     make(map[string]string),
		ProfessionalOptions: NewSchedulerProfessionalOptions(),
		Options:             NewSchedulerLitekubeOptions(),
		IgnoreOptions:       make(map[string]string),
	}
}

// delete keys already be disable or define in other block
func (opt *SchedulerOptions) CheckReservedOptions() error {
	// deep copy options
	optionsMap, oErr := common.StructToMap(opt.Options)
	if oErr != nil {
		return oErr
	}

	for k, _ := range optionsMap {
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

	for k, _ := range professionalOptionsMap {
		if value, ok := opt.ReservedOptions[k]; ok {
			opt.IgnoreOptions[k] = value
			delete(opt.ReservedOptions, k)
		}
	}
	return nil
}

func (opt *SchedulerOptions) HelpSection() *help.Section {
	section := help.NewSection("kube-scheduler", "kube-scheduler's startup parameters, we keep almost the same Settings as the original except logs relation: https://kubernetes.io/docs/reference/command-line-tools-reference/kube-scheduler/", nil)

	reserveSection := help.NewSection("reserve", "reserve parameters, you can still add args unmentioned before refer to kube-scheduler official website.", nil)
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

func (opt *SchedulerOptions) ToMap() (map[string]string, error) {
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
func (opt *SchedulerOptions) PrintFlags(prefix string, printFunc func(format string, a ...interface{}) error) error {
	// print flags
	flags, err := opt.ToMap()
	if err != nil {
		return err
	}
	printMap(flags, prefix, printFunc)
	// print flags-ignored
	printMap(opt.IgnoreOptions, fmt.Sprintf("%s-<%s>", prefix, UnreserveTip), printFunc)
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
