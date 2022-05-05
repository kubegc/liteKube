package runtime

import (
	"context"
	"fmt"

	// link to github.com/Litekube/kine, we have make some addition

	"github.com/litekube/LiteKube/pkg/options/worker/kubelet"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kubelet/app"
)

type Kubelet struct {
	ctx     context.Context
	LogPath string
	Options *kubelet.KubeletOptions
}

func NewKubelet(ctx context.Context, opt *kubelet.KubeletOptions, logPath string) *Kubelet {
	return &Kubelet{
		ctx:     ctx,
		Options: opt,
		LogPath: logPath,
	}
}

// start run in routine and no wait
func (s *Kubelet) Run() error {
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

	command := app.NewKubeletCommand(context.Background())
	command.SetArgs(argsValue)

	klog.Infof("==>kubelet: %s\n", argsValue)

	go func() {
		err := command.ExecuteContext(s.ctx)
		if err != nil {
			klog.Fatalf("kubelet exited: %v", err)
		}
	}()

	return nil
}
