package runtime

import (
	"context"
	"fmt"
	goruntime "runtime"

	// link to github.com/Litekube/kine, we have make some addition
	"github.com/litekube/LiteKube/pkg/logger"
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
	ptr, _, _, ok := goruntime.Caller(0)
	if ok {
		logger.DefaultLogger.SetLog(goruntime.FuncForPC(ptr).Name(), s.LogPath)
	} else {
		klog.Errorf("fail to init kine log")
	}

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

	fmt.Println("====>scheduler:", argsValue)

	go func() {
		err := command.Execute()
		if err != nil {
			fmt.Printf("kube-scheduler exited: %v", err)
			klog.Infof("kube-scheduler: %v", err)
			panic(err)
		}
	}()

	return nil
}
