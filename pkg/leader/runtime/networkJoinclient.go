package runtime

import (
	"context"

	"github.com/Litekube/network-controller/config"

	// link to github.com/Litekube/kine, we have make some addition
	"github.com/Litekube/network-controller/network"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
	"k8s.io/klog/v2"
)

type NetWorkJoinClient struct {
	ctx         context.Context
	WorkDir     string
	TLSDir      string
	LogDir      string
	BindAddress string
	Port        uint16
	CAPath      string
	CertPath    string
	KeyPath     string
	NodeToken   string
}

func NewNetWorkJoinClient(ctx context.Context, opt *netmanager.NetManagerOptions, workDir string, tlsDir string, logDir string) *NetWorkJoinClient {

	return &NetWorkJoinClient{
		ctx:         ctx,
		WorkDir:     workDir,
		TLSDir:      tlsDir,
		LogDir:      logDir,
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
	klog.Info("run network-controller client")

	client := network.NewClient(config.ClientConfig{
		CAFile:         s.CAPath,
		ClientCertFile: s.CertPath,
		ClientKeyFile:  s.KeyPath,

		WorkDir: s.WorkDir,
		LogDir:  s.LogDir,

		ServerAddr:      s.BindAddress,
		Port:            int(s.Port),
		MTU:             1400,
		Token:           s.NodeToken,
		RedirectGateway: false,
	})

	go func() {
		err := client.Run()
		if err != nil {
			klog.Infof("network-controller client exited: %v", err)
			panic(err)
		}

	}()

	return nil
}
