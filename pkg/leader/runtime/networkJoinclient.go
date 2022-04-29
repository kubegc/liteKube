package runtime

import (
	"context"
	goruntime "runtime"

	"github.com/Litekube/network-controller/config"
	// link to github.com/Litekube/kine, we have make some addition
	"github.com/Litekube/network-controller/network"
	"github.com/litekube/LiteKube/pkg/logger"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
	"k8s.io/klog/v2"
)

type NetWorkJoinClient struct {
	ctx         context.Context
	LogPath     string
	BindAddress string
	Port        uint16
	CAPath      string
	CertPath    string
	KeyPath     string
	NodeToken   string
	NetClient   config.ClientConfig
}

func NewNetWorkJoinClient(ctx context.Context, opt *netmanager.NetManagerOptions, logPath string) *NetWorkJoinClient {

	return &NetWorkJoinClient{
		ctx:         ctx,
		LogPath:     logPath,
		BindAddress: opt.JoinOptions.Address,
		Port:        opt.JoinOptions.SecurePort,
		CAPath:      opt.JoinOptions.CACert,
		CertPath:    opt.JoinOptions.ClientCertFile,
		KeyPath:     opt.JoinOptions.ClientkeyFile,
		NodeToken:   opt.NodeToken,
		NetClient: config.ClientConfig{
			CAFile:          opt.JoinOptions.CACert,
			ClientCertFile:  opt.JoinOptions.ClientCertFile,
			ClientKeyFile:   opt.JoinOptions.ClientkeyFile,
			ServerAddr:      opt.JoinOptions.Address,
			Port:            int(opt.JoinOptions.SecurePort),
			MTU:             1400,
			Token:           opt.NodeToken,
			RedirectGateway: false,
		},
	}

}

// start run in routine and no wait
func (s *NetWorkJoinClient) Run() error {
	ptr, _, _, ok := goruntime.Caller(0)
	if ok {
		logger.DefaultLogger.SetLog(goruntime.FuncForPC(ptr).Name(), s.LogPath)
	} else {
		klog.Errorf("fail to init kine log")
	}

	klog.Info("run network manager client")

	client := network.NewClient(s.NetClient)
	err := client.Run()
	if err != nil {
		return err
	}

	return nil
}
