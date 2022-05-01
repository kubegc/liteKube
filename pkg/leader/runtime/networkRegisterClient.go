package runtime

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Litekube/network-controller/grpc/grpc_client"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
)

//var NRClient *NetWorkRegisterClient = nil

type NetWorkRegisterClient struct {
	ctx             context.Context
	BindAddress     string
	Port            uint16
	CAPath          string
	CertPath        string
	KeyPath         string
	NodeToken       string
	NCClient        *grpc_client.GrpcClient
	BootstrapClient *grpc_client.GrpcBootStrapClient
}

func NewNetWorkRegisterClient(ctx context.Context, opt *netmanager.NetManagerOptions) *NetWorkRegisterClient {
	client := &NetWorkRegisterClient{
		ctx:         ctx,
		BindAddress: opt.RegisterOptions.Address,
		Port:        opt.RegisterOptions.SecurePort,
		CAPath:      opt.RegisterOptions.CACert,
		CertPath:    opt.RegisterOptions.ClientCertFile,
		KeyPath:     opt.RegisterOptions.ClientkeyFile,
		NodeToken:   opt.NodeToken,
		NCClient:    nil,
		BootstrapClient: &grpc_client.GrpcBootStrapClient{
			Ip:            "",
			BootstrapPort: "",
		},
	}

	client.NCClient = &grpc_client.GrpcClient{
		Ip:       client.BindAddress,
		Port:     strconv.FormatUint(uint64(client.Port), 10),
		CAFile:   client.CAPath,
		CertFile: client.CertPath,
		KeyFile:  client.KeyPath,
	}

	if err := client.NCClient.InitGrpcClientConn(); err != nil {
		panic(err)
	}

	// start in 5s
	for i := 0; i < 10; i++ {
		resp, err := client.NCClient.C.HealthCheck(ctx, &pb_gen.HealthCheckRequest{})
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if resp.Code == "200" {
			break
		}
		if i == 9 {
			panic(err)
		}
	}

	return client
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

	// // wait to be delete and implement functionality in network-controller server
	// if nodeToken == global.ReservedNodeToken {
	// 	return "10.1.1.1", nil
	// }

	if c.NCClient.C == nil {
		if err := c.NCClient.InitGrpcClientConn(); err != nil {
			panic(err)
		}
	}

	// start in 5s
	for i := 0; i < 10; i++ {
		resp, err := c.NCClient.C.HealthCheck(c.ctx, &pb_gen.HealthCheckRequest{})
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if resp.Code == "200" {
			break
		}
		if i == 9 {
			panic(err)
		}
	}

	req := &pb_gen.CheckConnStateRequest{
		Token: nodeToken,
	}

	resp, err := c.NCClient.C.CheckConnState(c.ctx, req)
	if err != nil {
		return "", err
	}

	if resp.BindIp == "" {
		return "", fmt.Errorf("none remote ip get")
	}

	return resp.BindIp, nil
}

func (c *NetWorkRegisterClient) CreateBootStrapToken(life int64) (string, error) {
	if c == nil {
		return "", fmt.Errorf("nil for NetWorkRegisterClient")
	}

	if c.NCClient.C == nil {
		if err := c.NCClient.InitGrpcClientConn(); err != nil {
			panic(err)
		}
	}

	// start in 5s
	for i := 0; i < 10; i++ {
		resp, err := c.NCClient.C.HealthCheck(c.ctx, &pb_gen.HealthCheckRequest{})
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if resp.Code == "200" {
			break
		}
		if i == 9 {
			panic(err)
		}
	}

	req := &pb_gen.GetBootStrapTokenRequest{
		ExpireTime: life,
	}

	resp, err := c.NCClient.C.GetBootStrapToken(c.ctx, req)
	if err != nil {
		return "", err
	}

	c.BootstrapClient.BootstrapPort = resp.Port
	c.BootstrapClient.Ip = resp.CloudIp
	err = c.BootstrapClient.InitGrpcBootstrapClientConn()
	if err != nil {
		return "", err
	}

	return resp.BootStrapToken, nil
}

func (c *NetWorkRegisterClient) GetBootStrapAddress() (string, error) {
	if c == nil {
		return "", fmt.Errorf("nil for NetWorkRegisterClient")
	}

	if c.BootstrapClient.Ip == "" {
		if _, err := c.CreateBootStrapToken(0); err != nil {
			return "", err
		}
	}

	return c.BootstrapClient.Ip, nil
}

func (c *NetWorkRegisterClient) GetBootStrapPort() (uint16, error) {
	if c == nil {
		return 0, fmt.Errorf("nil for NetWorkRegisterClient")
	}

	if c.BootstrapClient.BootstrapPort == "" {
		if _, err := c.CreateBootStrapToken(0); err != nil {
			return 0, err
		}
	}

	port, _ := strconv.ParseUint(c.BootstrapClient.BootstrapPort, 10, 16)
	return uint16(port), nil
}
