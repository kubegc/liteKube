package config

import options "github.com/litekube/LiteKube/pkg/options/leader"

type LeaderRuntime struct {
	options *options.LeaderOptions
}

func NewLeaderRuntime() *LeaderRuntime {
	return &LeaderRuntime{
		options: options.NewLeaderOptions(),
	}
}

func (runtime *LeaderRuntime) LoadConfig(opt *options.LeaderOptions) {
	opt.
		runtime.options = opt
}
