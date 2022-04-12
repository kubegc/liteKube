package leader

import "github.com/spf13/pflag"

var ConfigFile string

const ConfigFileFlagName = "config-file"

func AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&ConfigFile, ConfigFileFlagName, "", "YAML File to store leader startup parameters")
}
