package runtime

import (
	"context"
	"fmt"
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

	fmt.Printf("node-token here %s\n", s.NodeToken)
	client := network.NewClient(config.ClientConfig{
		CAFile:          s.CAPath,
		ClientCertFile:  s.CertPath,
		ClientKeyFile:   s.KeyPath,
		ServerAddr:      s.BindAddress,
		Port:            int(s.Port),
		MTU:             1400,
		Token:           s.NodeToken,
		RedirectGateway: false,
	})

	go func() {
		err := client.Run()
		if err != nil {
			fmt.Printf("network controller client exited: %v", err)
			klog.Infof("network controller client exited: %v", err)
			panic(err)
		}
	}()

	return nil
}
