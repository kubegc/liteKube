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
	NodeToken   string
}

func NewNetWorkRegisterClient(ctx context.Context, opt *netmanager.NetManagerOptions) *NetWorkRegisterClient {
	NRClient = &NetWorkRegisterClient{
		ctx:         ctx,
		BindAddress: opt.RegisterOptions.Address,
		Port:        opt.RegisterOptions.SecurePort,
		CAPath:      opt.RegisterOptions.CACert,
		CertPath:    opt.RegisterOptions.ClientCertFile,
		KeyPath:     opt.RegisterOptions.ClientkeyFile,
		NodeToken:   opt.NodeToken,
	}

	return NRClient
}

// query local ip
func (c *NetWorkRegisterClient) QueryIp() (string, error) {
	return c.QueryIpByToken(c.NodeToken)
}

// query ip by node-token
func (c *NetWorkRegisterClient) QueryIpByToken(nodeToken string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("nil for NetWorkRegisterClient")
	}

	// grpc remote
	return "192.168.154.101", nil
}

func (c *NetWorkRegisterClient) CreateBootStrapToken(life int64) (string, error) {
	if c == nil {
		return "", fmt.Errorf("nil for NetWorkRegisterClient")
	}

	// grpc remote
	return "this-is-test-token", nil
}

func (c *NetWorkRegisterClient) GetBootStrapAddress() (string, error) {
	if c == nil {
		return "", fmt.Errorf("nil for NetWorkRegisterClient")
	}

	// grpc remote
	return c.BindAddress, nil
}

func (c *NetWorkRegisterClient) GetBootStrapPort() (uint16, error) {
	if c == nil {
		return 0, fmt.Errorf("nil for NetWorkRegisterClient")
	}

	// grpc remote
	return c.Port, nil
}
