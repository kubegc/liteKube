package config

import options "github.com/litekube/LiteKube/pkg/options/leader"

type LeaderRuntime struct {
	options *options.LeaderOptions
}

func NewLeaderRuntime() *LeaderRuntime {
	return &LeaderRuntime{
		options: nil,
	}
}

func (runtime *LeaderRuntime) LoadConfig(opt *options.LeaderOptions) {
	runtime.options = opt
}
