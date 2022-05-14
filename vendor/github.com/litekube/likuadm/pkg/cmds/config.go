package cmds

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/leader/runtime/control"
	"github.com/litekube/likuadm/pkg/authentication"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	Client control.LeaderControlClient
}

func NewClient() *Client {
	client := &Client{}

	if err := client.Init(); err != nil {
		panic(err)
	}

	return client
}

func (c *Client) Init() error {
	if c == nil {
		return nil
	}

	auth := &authentication.TokenAuthentication{
		Token: GlobalConfig.token,
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", GlobalConfig.ip, GlobalConfig.port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(auth))
	if err != nil {
		return err
	}

	c.Client = control.NewLeaderControlClient(conn)
	c.conn = conn
	return nil
}

func (c *Client) GRPC() control.LeaderControlClient {
	return c.Client
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
