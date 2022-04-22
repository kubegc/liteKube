package runtime

import (
	"context"
	"net/http"
	"path/filepath"

	// link to github.com/Litekube/kine, we have make some addition

	"github.com/litekube/LiteKube/pkg/options/leader/apiserver"
	"github.com/litekube/LiteKube/pkg/options/leader/controllermanager"
	"github.com/litekube/LiteKube/pkg/options/leader/scheduler"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/klog/v2"
)

type KubernatesServer struct {
	ctx               context.Context
	logPath           string
	apiserver         *Apiserver
	controllerManager *ControllerManager
	scheduler         *Scheduler
	ApiserverOptions  *apiserver.ApiserverOptions
	ControllerOptions *controllermanager.ControllerManagerOptions
	SchedulerOptions  *scheduler.SchedulerOptions
	KubeAdminPath     string
}

func NewKubernatesServer(ctx context.Context, apiserverOptions *apiserver.ApiserverOptions, controllerOptions *controllermanager.ControllerManagerOptions, schedulerOptions *scheduler.SchedulerOptions, kubeAdminPath string, logPath string) *KubernatesServer {
	return &KubernatesServer{
		ctx:               ctx,
		logPath:           logPath,
		apiserver:         NewApiserver(ctx, apiserverOptions, filepath.Join(logPath, "kube-apiserver.log")),
		controllerManager: NewControllerManager(ctx, controllerOptions, filepath.Join(logPath, "kube-controller-manager.log")),
		scheduler:         NewScheduler(ctx, schedulerOptions, filepath.Join(logPath, "kube-scheduler.log")),
		ApiserverOptions:  apiserverOptions,
		ControllerOptions: controllerOptions,
		SchedulerOptions:  schedulerOptions,
		KubeAdminPath:     kubeAdminPath,
	}
}

// start run in routine and no wait
func (s *KubernatesServer) Run() error {
	klog.Info("start to run kubernates server")

	if err := s.apiserver.Run(); err != nil {
		return err
	}

	if err := s.controllerManager.Run(s.KubeAdminPath); err != nil {
		return err
	}

	if err := s.scheduler.Run(); err != nil {
		return err
	}

	return nil
}

func (s *KubernatesServer) StartUpConfig() (*http.Handler, *authenticator.Request) {
	return s.apiserver.StartUpConfig()
}
