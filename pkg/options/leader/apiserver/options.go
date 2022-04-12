package apiserver

import "github.com/litekube/LiteKube/pkg/help"

// struct to store args from input
type ApiserverOptions struct {
	ReservedOptions     map[string]string             `yaml:"reserve"`
	ProfessionalOptions *ApiserverProfessionalOptions `yaml:"professional"`
	Options             *ApiserverLitekubeOptions     `yaml:"options"`
}

func NewApiserverOptions() *ApiserverOptions {
	return &ApiserverOptions{
		ReservedOptions:     make(map[string]string),
		ProfessionalOptions: NewApiserverProfessionalOptions(),
		Options:             NewApiserverLitekubeOptions(),
	}
}

// delete keys already be disable or define in other block
// return all the bad keys, nil
func (opt *ApiserverOptions) CheckReservedOptions(banArgs []string) ([]string, error) {
	checkArgs := append(UnreservedArgs, banArgs...)

	args := make([]string, 0, len(checkArgs))
	for _, arg := range checkArgs {
		if _, ok := opt.ReservedOptions[arg]; ok {
			args = append(args, arg)
			delete(opt.ReservedOptions, arg)
		}
	}

	return args, nil
}

func (opt *ApiserverOptions) HelpSection() *help.Section {
	section := help.NewSection("kube-apiserver", "kube-Apiserver's startup parameters, we keep almost the same Settings as the original except logs relation: https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/", nil)

	reserveSection := help.NewSection("reserve", "reserve parameters, you can still add args unmentioned before refer to kube-apiserver official website.", nil)
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
