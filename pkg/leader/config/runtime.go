package config

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/litekube/LiteKube/pkg/leader/runtime"
	options "github.com/litekube/LiteKube/pkg/options/leader"
	"k8s.io/klog/v2"
)

type LeaderRuntime struct {
	control                 *ControlSignal
	FlagsOption             *options.LeaderOptions
	RuntimeOption           *options.LeaderOptions
	RuntimeAuthentication   *RuntimeAuthentications
	KineServer              *runtime.KineServer
	NetworkControllerServer *runtime.NetWorkControllerServer
	NetworkJoinClient       *runtime.NetWorkJoinClient
	NetworkRegisterClient   *runtime.NetWorkRegisterClient
	KubernetesServer        *runtime.KubernatesServer
	controlServer           *runtime.LiteKubeControl
	OwnKineCert             bool
}

// control progress end
type ControlSignal struct {
	ctx  context.Context
	stop context.CancelFunc
	wg   *sync.WaitGroup
}

func NewLeaderRuntime(flags *options.LeaderOptions) *LeaderRuntime {
	ctx, stop := context.WithCancel(context.TODO())
	return &LeaderRuntime{
		control: &ControlSignal{
			ctx:  ctx,
			stop: stop,
			wg:   &sync.WaitGroup{},
		},
		FlagsOption:             flags,
		RuntimeOption:           options.NewLeaderOptions(),
		RuntimeAuthentication:   nil,
		NetworkControllerServer: nil,
		NetworkJoinClient:       nil,
		NetworkRegisterClient:   nil,
		KubernetesServer:        nil,
		OwnKineCert:             false,
		controlServer:           nil,
		KineServer:              nil,
	}
}

// run kine server, network manager server, network client
func (leaderRuntime *LeaderRuntime) RunForward() error {
	defer leaderRuntime.Done()
	leaderRuntime.Add()

	if leaderRuntime.RuntimeOption.GlobalOptions.RunKine {
		// run kine and network-manager
		leaderRuntime.KineServer = runtime.NewKineServer(leaderRuntime.control.ctx,
			leaderRuntime.RuntimeOption.KineOptions,
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/kine/"),
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/kine.log"),
		)
		if err := leaderRuntime.KineServer.Run(); err != nil {
			klog.Errorf("bad args for kine server")
			return err
		}
	}

	if leaderRuntime.RuntimeOption.GlobalOptions.RunNetManager {
		leaderRuntime.NetworkControllerServer = runtime.NewNetWorkControllerServer(leaderRuntime.control.ctx,
			leaderRuntime.RuntimeAuthentication.NetWorkController,
			leaderRuntime.RuntimeOption.NetmamagerOptions,
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/network-manager.log"),
		)
		if err := leaderRuntime.NetworkControllerServer.Run(); err != nil {
			klog.Errorf("bad args for network manager server")
			return err
		}
	} else {
		leaderRuntime.NetworkJoinClient = runtime.NewNetWorkJoinClient(leaderRuntime.control.ctx,
			leaderRuntime.RuntimeOption.NetmamagerOptions,
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/network-client.log"),
		)
		if err := leaderRuntime.NetworkJoinClient.Run(); err != nil {
			klog.Errorf("bad args for network manager client")
			return err
		}
	}

	// wait to be enhance by network-controller
	time.Sleep(10 * time.Second) // only for debug, waiting for network-controller to start

	leaderRuntime.NetworkRegisterClient = runtime.NewNetWorkRegisterClient(leaderRuntime.control.ctx, leaderRuntime.RuntimeOption.NetmamagerOptions)
	return nil
}

// run k8s and litekube controller
func (leaderRuntime *LeaderRuntime) Run() error {
	defer leaderRuntime.Done()
	leaderRuntime.Add()

	leaderRuntime.KubernetesServer = runtime.NewKubernatesServer(leaderRuntime.control.ctx,
		leaderRuntime.RuntimeOption.ApiserverOptions,
		leaderRuntime.RuntimeOption.ControllerManagerOptions,
		leaderRuntime.RuntimeOption.SchedulerOptions,
		leaderRuntime.RuntimeAuthentication.Kubernetes.KubeConfigAdmin,
		filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/kubernetes/"),
	)

	if err := leaderRuntime.KubernetesServer.Run(); err != nil {
		klog.Errorf("fail to start kubernetes server. Error: %s", err.Error())
		return err
	}

	leaderRuntime.controlServer = runtime.NewLiteKubeControl(leaderRuntime.control.ctx,
		leaderRuntime.NetworkRegisterClient,
		filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "tls/buffer"),
		leaderRuntime.RuntimeOption.NetmamagerOptions.NodeToken,
		fmt.Sprintf("https://%s:%d", leaderRuntime.RuntimeOption.ApiserverOptions.ProfessionalOptions.AdvertiseAddress, leaderRuntime.RuntimeOption.ApiserverOptions.Options.SecurePort),
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.RootCaFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeApiserverClientCertFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeletClientCertFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeletClientKeyFile,
		leaderRuntime.RuntimeOption.ApiserverOptions.ProfessionalOptions.TokenAuthFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.Options.ClusterCidr,
	)

	if err := leaderRuntime.controlServer.Run(); err != nil {
		klog.Errorf("fail to start litekube control server. Error: %s", err.Error())
		return err
	}

	return nil
}

func (leaderRuntime *LeaderRuntime) Stop() error {
	defer leaderRuntime.Wait()

	// give signal to end process
	leaderRuntime.control.stop()

	// stop while all return
	return nil
}

func (leaderRuntime *LeaderRuntime) Done() {
	leaderRuntime.control.wg.Done()
}

func (leaderRuntime *LeaderRuntime) Wait() {
	leaderRuntime.control.wg.Wait()
}

func (leaderRuntime *LeaderRuntime) Add() {
	leaderRuntime.control.wg.Add(1)
}
