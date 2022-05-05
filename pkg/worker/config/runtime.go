package config

import (
	"context"
	"path/filepath"
	"sync"

	leaderruntime "github.com/litekube/LiteKube/pkg/leader/runtime"
	options "github.com/litekube/LiteKube/pkg/options/worker"
	workerruntime "github.com/litekube/LiteKube/pkg/worker/runtime"
	"k8s.io/klog/v2"
)

type WorkerRuntime struct {
	control               *ControlSignal
	FlagsOption           *options.WorkerOptions
	RuntimeOption         *options.WorkerOptions
	RuntimeAuthentication *RuntimeAuthentications
	NetworkJoinClient     *leaderruntime.NetWorkJoinClient
	NetworkRegisterClient *leaderruntime.NetWorkRegisterClient
	KubernatesClient      *workerruntime.KubernatesClient
}

// control progress end
type ControlSignal struct {
	ctx  context.Context
	stop context.CancelFunc
	wg   *sync.WaitGroup
}

func NewWorkerRuntime(flags *options.WorkerOptions) *WorkerRuntime {
	ctx, stop := context.WithCancel(context.TODO())
	return &WorkerRuntime{
		control: &ControlSignal{
			ctx:  ctx,
			stop: stop,
			wg:   &sync.WaitGroup{},
		},
		FlagsOption:           flags,
		RuntimeOption:         options.NewWorkerOptions(),
		RuntimeAuthentication: nil,
		NetworkJoinClient:     nil,
		NetworkRegisterClient: nil,
	}
}

// run kine server, network manager server, network client
func (workerRuntime *WorkerRuntime) RunForward() error {
	defer workerRuntime.Done()
	workerRuntime.Add()

	if workerRuntime.RuntimeOption.NetmamagerOptions.Token != "local" {
		workerRuntime.NetworkJoinClient = leaderruntime.NewNetWorkJoinClient(workerRuntime.control.ctx,
			workerRuntime.RuntimeOption.NetmamagerOptions,
			filepath.Join(workerRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/network-client.log"),
		)

		if err := workerRuntime.NetworkJoinClient.Run(); err != nil {
			klog.Errorf("bad args for network manager client")
			return err
		}
	}

	workerRuntime.NetworkRegisterClient = leaderruntime.NewNetWorkRegisterClient(workerRuntime.control.ctx, workerRuntime.RuntimeOption.NetmamagerOptions)
	return nil
}

// run k8s
func (workerRuntime *WorkerRuntime) Run() error {
	defer workerRuntime.Done()
	workerRuntime.Add()

	workerRuntime.KubernatesClient = workerruntime.NewKubernatesClient(workerRuntime.control.ctx,
		workerRuntime.RuntimeOption.KubeletOptions,
		workerRuntime.RuntimeOption.KubeProxyOptions,
		filepath.Join(workerRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/kubernetes/"),
	)

	if err := workerRuntime.KubernatesClient.Run(); err != nil {
		klog.Fatal("fail to start kubernetes node. Error: %s", err.Error())
		return err
	}

	return nil
}

func (workerRuntime *WorkerRuntime) Stop() error {
	defer workerRuntime.Wait()

	// give signal to end process
	workerRuntime.control.stop()

	// stop while all return
	return nil
}

func (workerRuntime *WorkerRuntime) Done() {
	workerRuntime.control.wg.Done()
}

func (workerRuntime *WorkerRuntime) Wait() {
	workerRuntime.control.wg.Wait()
}

func (workerRuntime *WorkerRuntime) Add() {
	workerRuntime.control.wg.Add(1)
}
