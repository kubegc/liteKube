package runtime

import (
	"context"
	"fmt"

	// link to github.com/Litekube/kine, we have make some addition

	"github.com/litekube/LiteKube/pkg/options/leader/scheduler"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

type Scheduler struct {
	ctx     context.Context
	LogPath string
	Options *scheduler.SchedulerOptions
}

func NewScheduler(ctx context.Context, opt *scheduler.SchedulerOptions, logPath string) *Scheduler {
	return &Scheduler{
		ctx:     ctx,
		Options: opt,
		LogPath: logPath,
	}
}

// start run in routine and no wait
func (s *Scheduler) Run() error {
	klog.Info("run kube-scheduler:")

	args, err := s.Options.ToMap()
	if err != nil {
		return err
	}

	argsValue := make([]string, 0, len(args))
	for k, v := range args {
		if v == "-" || v == "" {
			continue
		}
		argsValue = append(argsValue, fmt.Sprintf("--%s=%s", k, v))
	}

	command := app.NewSchedulerCommand()
	command.SetArgs(argsValue)

	klog.Infof("==>kube-scheduler: %s\n", argsValue)

	go func() {
		err := command.ExecuteContext(s.ctx)
		if err != nil {
			fmt.Printf("kube-scheduler exited: %v", err)
			klog.Infof("kube-scheduler: %v", err)
			panic(err)
		}
	}()

	return nil
}
