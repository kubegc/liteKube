package config

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/litekube/LiteKube/pkg/leader/runtime"
	leaderruntime "github.com/litekube/LiteKube/pkg/leader/runtime"
	options "github.com/litekube/LiteKube/pkg/options/worker"
	"k8s.io/klog/v2"
)

type WorkerRuntime struct {
	control               *ControlSignal
	FlagsOption           *options.WorkerOptions
	RuntimeOption         *options.WorkerOptions
	RuntimeAuthentication *RuntimeAuthentications
	NetworkJoinClient     *leaderruntime.NetWorkJoinClient
	NetworkRegisterClient *leaderruntime.NetWorkRegisterClient
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
		workerRuntime.NetworkJoinClient = runtime.NewNetWorkJoinClient(workerRuntime.control.ctx,
			workerRuntime.RuntimeOption.NetmamagerOptions,
			filepath.Join(workerRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/network-client.log"),
		)

		if err := workerRuntime.NetworkJoinClient.Run(); err != nil {
			klog.Errorf("bad args for network manager client")
			return err
		}
	}

	workerRuntime.NetworkRegisterClient = runtime.NewNetWorkRegisterClient(workerRuntime.control.ctx, workerRuntime.RuntimeOption.NetmamagerOptions)
	return nil
}

// run k8s
func (workerRuntime *WorkerRuntime) Run() error {
	defer workerRuntime.Done()
	workerRuntime.Add()

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
