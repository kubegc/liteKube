package config

import (
	options "github.com/litekube/LiteKube/pkg/options/leader"
)

type LeaderRuntime struct {
	FlagsOption           *options.LeaderOptions
	RuntimeOption         *options.LeaderOptions
	RuntimeAuthentication *RuntimeAuthentications
	OwnKineCert           bool
}

type RuntimeOptions struct {
	*options.LeaderOptions
	OwnKineCert bool
}

func NewLeaderRuntime(flags *options.LeaderOptions) *LeaderRuntime {
	return &LeaderRuntime{
		FlagsOption:           flags,
		RuntimeOption:         options.NewLeaderOptions(),
		RuntimeAuthentication: nil,
		OwnKineCert:           false,
	}
}

// func NewRuntimeOptions() *RuntimeOptions {
// 	return &RuntimeOptions{
// 		LeaderOptions: options.NewLeaderOptions(),

// 		//Logger:        nil,
// 	}
// }

func (runtime *LeaderRuntime) CheckArgs() {

}

// // check kine args
// func (runtime *LeaderRuntime) CheckKine() error {
// 	// disable kine
// 	if !runtime.FlagsOption.GlobalOptions.RunKine {
// 		runtime.FlagsOption.KineOptions = nil
// 		return nil
// 	}

// 	// enable kine
// 	return nil
// }
