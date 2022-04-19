package runtime

import (
	"context"

	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
)

type NetWorkRegisterClient struct {
	ctx         context.Context
	BindAddress string
	Port        uint16
	CAPath      string
	CertPath    string
	KeyPath     string
}

func NewNetWorkRegisterClient(ctx context.Context, opt *netmanager.NetManagerOptions) *NetWorkRegisterClient {
	return &NetWorkRegisterClient{
		ctx:         ctx,
		BindAddress: opt.RegisterOptions.Address,
		Port:        opt.RegisterOptions.SecurePort,
		CAPath:      opt.RegisterOptions.CACert,
		CertPath:    opt.RegisterOptions.ClientCertFile,
		KeyPath:     opt.RegisterOptions.ClientkeyFile,
	}
}

func (c *NetWorkRegisterClient) QueryIp() (string, error) {
	// grpc remote
	return "192.168.154.101", nil
}
