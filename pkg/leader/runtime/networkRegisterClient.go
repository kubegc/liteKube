package runtime

import (
	"context"
	"fmt"

	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
)

var NRClient *NetWorkRegisterClient = nil

type NetWorkRegisterClient struct {
	ctx         context.Context
	BindAddress string
	Port        uint16
	CAPath      string
	CertPath    string
	KeyPath     string
}

func NewNetWorkRegisterClient(ctx context.Context, opt *netmanager.NetManagerOptions) *NetWorkRegisterClient {
	NRClient = &NetWorkRegisterClient{
		ctx:         ctx,
		BindAddress: opt.RegisterOptions.Address,
		Port:        opt.RegisterOptions.SecurePort,
		CAPath:      opt.RegisterOptions.CACert,
		CertPath:    opt.RegisterOptions.ClientCertFile,
		KeyPath:     opt.RegisterOptions.ClientkeyFile,
	}

	return NRClient
}

func (c *NetWorkRegisterClient) QueryIp() (string, error) {
	if c == nil {
		return "", fmt.Errorf("nil for NetWorkRegisterClient")
	}

	// grpc remote
	return "192.168.154.101", nil
}
