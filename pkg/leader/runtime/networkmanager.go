package runtime

import (
	"context"
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/network"
	goruntime "runtime"
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

	NCServer config.ServerConfig
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

		NCServer: config.ServerConfig{
			Ip:   opt.RegisterBindAddress,
			Port: int(clientOpt.JoinOptions.SecurePort),
			// todo config BootstrapPort
			BootstrapPort: 6439,
			GrpcPort:      int(clientOpt.RegisterOptions.SecurePort),
			// todo config NetworkAddr
			NetworkAddr:     "10.1.1.1/24",
			MTU:             1400,
			Interconnection: false,

			NetworkCAFile:         opt.JoinCACert,
			NetworkCAKeyFile:      opt.JoinCAKey,
			NetworkServerCertFile: opt.JoinServerCert,
			NetworkServerKeyFile:  opt.JoinServerkey,

			GrpcCAFile:         opt.RegisterCACert,
			GrpcCAKeyFile:      opt.RegisterCAKey,
			GrpcServerCertFile: opt.RegisterServerCert,
			GrpcServerKeyFile:  opt.RegisterServerkey,
		},
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

	server := network.NewServer(s.NCServer)
	err := server.Run()
	if err != nil {
		return err
	}

	return nil
}
