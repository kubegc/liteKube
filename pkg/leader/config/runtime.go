package config

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	workerapp "github.com/litekube/LiteKube/cmd/worker/app"
	"github.com/litekube/LiteKube/pkg/leader/runtime"
	options "github.com/litekube/LiteKube/pkg/options/leader"
	workeroptions "github.com/litekube/LiteKube/pkg/options/worker"
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
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "kine/"),
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "logs/kine.log"),
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
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "network-controller/server/"),
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "tls/network-control/"),
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "logs/network-controller/server/"),
		)
		if err := leaderRuntime.NetworkControllerServer.Run(); err != nil {
			klog.Errorf("bad args for network-controller server")
			return err
		}
	} else {
		leaderRuntime.NetworkJoinClient = runtime.NewNetWorkJoinClient(leaderRuntime.control.ctx,
			leaderRuntime.RuntimeOption.NetmamagerOptions,
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "network-controller/client/"),
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "tls/network-control/"),
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "logs/network-controller/client/"),
		)
		if err := leaderRuntime.NetworkJoinClient.Run(); err != nil {
			klog.Errorf("bad args for network-controller client")
			return err
		}
	}

	// wait to be enhance by network-controller
	//time.Sleep(10 * time.Second) // only for debug, waiting for network-controller to start

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
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeletServingCertFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeletServingKeyFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeletClientCertFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeletClientKeyFile,
		leaderRuntime.RuntimeOption.ApiserverOptions.ProfessionalOptions.TokenAuthFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.Options.ClusterCidr,
		leaderRuntime.RuntimeOption.ApiserverOptions.Options.ServiceClusterIpRange,
	)

	if err := leaderRuntime.controlServer.Run(); err != nil {
		klog.Fatal("fail to start litekube control server. Error: %s", err.Error())
		return err
	}

	// run worker
	if leaderRuntime.RuntimeOption.GlobalOptions.EnableWorker {
		workeroptions.ConfigFile = leaderRuntime.RuntimeOption.GlobalOptions.WorkerConfig
		workerOpt := workeroptions.NewWorkerOptions()
		if err := workerOpt.LoadConfig(); err != nil {
			klog.Errorf("fail to run worker")
			return err
		}

		workerOpt.GlobalOptions.WorkDir = leaderRuntime.RuntimeOption.GlobalOptions.WorkDir
		workerOpt.GlobalOptions.LogDir = leaderRuntime.RuntimeOption.GlobalOptions.LogDir
		workerOpt.GlobalOptions.LogToStd = leaderRuntime.RuntimeOption.GlobalOptions.LogToStd
		workerOpt.GlobalOptions.LogToDir = leaderRuntime.RuntimeOption.GlobalOptions.LogToDir
		workerOpt.NetmamagerOptions.Token = "local"
		workerOpt.NetmamagerOptions.NodeToken = leaderRuntime.RuntimeOption.NetmamagerOptions.NodeToken
		workerOpt.NetmamagerOptions.RegisterOptions.Address = leaderRuntime.RuntimeOption.NetmamagerOptions.RegisterOptions.Address
		workerOpt.NetmamagerOptions.RegisterOptions.SecurePort = leaderRuntime.RuntimeOption.NetmamagerOptions.RegisterOptions.SecurePort
		workerOpt.NetmamagerOptions.RegisterOptions.CACert = leaderRuntime.RuntimeOption.NetmamagerOptions.RegisterOptions.CACert
		workerOpt.NetmamagerOptions.RegisterOptions.ClientCertFile = leaderRuntime.RuntimeOption.NetmamagerOptions.RegisterOptions.ClientCertFile
		workerOpt.NetmamagerOptions.RegisterOptions.ClientkeyFile = leaderRuntime.RuntimeOption.NetmamagerOptions.RegisterOptions.ClientkeyFile
		workerOpt.NetmamagerOptions.JoinOptions.Address = leaderRuntime.RuntimeOption.NetmamagerOptions.JoinOptions.Address
		workerOpt.NetmamagerOptions.JoinOptions.SecurePort = leaderRuntime.RuntimeOption.NetmamagerOptions.JoinOptions.SecurePort
		workerOpt.NetmamagerOptions.JoinOptions.CACert = leaderRuntime.RuntimeOption.NetmamagerOptions.JoinOptions.CACert
		workerOpt.NetmamagerOptions.JoinOptions.ClientCertFile = leaderRuntime.RuntimeOption.NetmamagerOptions.JoinOptions.ClientCertFile
		workerOpt.NetmamagerOptions.JoinOptions.ClientkeyFile = leaderRuntime.RuntimeOption.NetmamagerOptions.JoinOptions.ClientkeyFile
		workerOpt.GlobalOptions.LeaderToken = fmt.Sprintf("%s@local", leaderRuntime.RuntimeOption.NetmamagerOptions.NodeToken)
		go func() {
			// go func() {
			// 	// max 12-hour to wait worker start unless restart leader
			// 	for i := 0; i < 8640; i++ {
			// 		kubeconfig, err := clientcmd.BuildConfigFromFlags("", leaderRuntime.RuntimeAuthentication.Kubernetes.KubeConfigAdmin)
			// 		if err != nil {
			// 			klog.Errorf("bad kubeconfig while listen node-csr for leader")
			// 			return
			// 		}

			// 		k8sClient, err := kubernetes.NewForConfig(kubeconfig)
			// 		if err != nil {
			// 			klog.Errorf("fail to create client while listen node-csr for leader")
			// 			return
			// 		}

			// 		if csrList, err := k8sClient.CertificatesV1().CertificateSigningRequests().List(leaderRuntime.control.ctx, metav1.ListOptions{}); err != nil {
			// 			klog.Errorf("fail to get csr-list")
			// 			return
			// 		} else {
			// 			for _, csr := range csrList.Items {
			// 				block, rest := pem.Decode(csr.Spec.Request)
			// 				if block == nil || len(rest) > 0 {
			// 					klog.Errorf("get bad csr request")
			// 					return
			// 				}

			// 				key, _ := x509.ParseCertificateRequest(block.Bytes)
			// 				current_ip, err := leaderRuntime.NetworkRegisterClient.QueryIp()
			// 				if err != nil {
			// 					klog.Errorf("fail to get current ip")
			// 					return
			// 				}
			// 				if strings.Contains(key.Subject.CommonName, current_ip) {
			// 					k8sClient.CertificatesV1().CertificateSigningRequests().up
			// 					// k8sClient.CertificatesV1().CertificateSigningRequests().ApplyStatus(leaderRuntime.control.ctx,)
			// 				}
			// 			}
			// 		}

			// 		time.Sleep(5 * time.Second)
			// 	}
			// }()

			err := workerapp.Run(workerOpt, leaderRuntime.control.ctx.Done())
			if err != nil {
				klog.Errorf("==> worker exit: %v", err)
			}
		}()

		// kubectl certificate approve node-csr-xxxxxx
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
