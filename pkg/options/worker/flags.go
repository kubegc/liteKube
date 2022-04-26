package worker

import "github.com/spf13/pflag"

var ConfigFile string

// var JoinToken string

const ConfigFileFlagName = "config-file"

func AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&ConfigFile, ConfigFileFlagName, "", "YAML File to store leader startup parameters")
	// fs.StringVar(&JoinToken, "join", "", "[only need for the first time]. string use to join to cluster, you can get from leader by seek administrator for help.")
}
