package runtime

import (
	"context"
	"fmt"
	goruntime "runtime"

	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/network"

	// link to github.com/Litekube/kine, we have make some addition
	"github.com/litekube/LiteKube/pkg/leader/authentication"
	"github.com/litekube/LiteKube/pkg/logger"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
	"k8s.io/klog/v2"
)

type NetWorkManager struct {
	ctx     context.Context
	LogPath string

	// register
	RegisterBindAddress string
	RegisterPort        uint16
	RegisterCACert      string
	RegisterCAKey       string
	RegisterServerCert  string
	RegisterServerkey   string

	// join
	JoinBindAddress string
	JoinPort        uint16
	JoinCACert      string
	JoinCAKey       string
	JoinServerCert  string
	JoinServerkey   string
}

func NewNetWorkManager(ctx context.Context, opt *authentication.NetworkAuthentication, clientOpt *netmanager.NetManagerOptions, logPath string) *NetWorkManager {

	return &NetWorkManager{
		ctx:     ctx,
		LogPath: logPath,

		// register
		RegisterBindAddress: opt.RegisterBindAddress,
		RegisterPort:        clientOpt.RegisterOptions.SecurePort,
		RegisterCACert:      opt.RegisterCACert,
		RegisterCAKey:       opt.RegisterCAKey,
		RegisterServerCert:  opt.RegisterServerCert,
		RegisterServerkey:   opt.RegisterServerkey,

		// join
		JoinBindAddress: opt.JoinBindAddress,
		JoinPort:        clientOpt.JoinOptions.SecurePort,
		JoinCACert:      opt.JoinCACert,
		JoinCAKey:       opt.JoinCAKey,
		JoinServerCert:  opt.JoinServerCert,
		JoinServerkey:   opt.JoinServerkey,
	}
}

// start run in routine and no wait
func (s *NetWorkManager) Run() error {
	ptr, _, _, ok := goruntime.Caller(0)
	if ok {
		logger.DefaultLogger.SetLog(goruntime.FuncForPC(ptr).Name(), s.LogPath)
	} else {
		klog.Errorf("fail to init kine log")
	}

	klog.Info("run network manager")

	server := network.NewServer(config.ServerConfig{
		Ip:   s.RegisterBindAddress,
		Port: int(s.JoinPort),
		// todo config BootstrapPort
		BootstrapPort: 6439,
		GrpcPort:      int(s.RegisterPort),
		// todo config NetworkAddr
		NetworkAddr:     "10.1.1.1/24",
		MTU:             1400,
		Interconnection: false,

		NetworkCAFile:         s.JoinCACert,
		NetworkCAKeyFile:      s.JoinCAKey,
		NetworkServerCertFile: s.JoinServerCert,
		NetworkServerKeyFile:  s.JoinServerkey,

		GrpcCAFile:         s.RegisterServerCert,
		GrpcCAKeyFile:      s.RegisterCAKey,
		GrpcServerCertFile: s.RegisterServerCert,
		GrpcServerKeyFile:  s.RegisterServerkey,
	})

	go func() {
		err := server.Run()
		if err != nil {
			fmt.Printf("network controller exited: %v", err)
			klog.Infof("network controller exited: %v", err)
			panic(err)
		}

		s.ctx.Done()
	}()

	return nil
}
