package scheduler

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// options for Litekube to start kube-scheduler
type SchedulerLitekubeOptions struct {
	Profiling bool `yaml:"profiling"`
}

var defaultSLO SchedulerLitekubeOptions = SchedulerLitekubeOptions{
	Profiling: false,
}

func NewSchedulerLitekubeOptions() *SchedulerLitekubeOptions {
	options := defaultSLO
	return &options
}

func (opt *SchedulerLitekubeOptions) AddTips(section *help.Section) {

	section.AddTip("profiling", "bool", "deprecated. Enable profiling via web interface host:port/debug/pprof/.", fmt.Sprintf("%t", defaultSLO.Profiling))
}
