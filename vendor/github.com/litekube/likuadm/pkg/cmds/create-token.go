package cmds

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/litekube/LiteKube/pkg/leader/runtime/control"
	"github.com/urfave/cli/v2"
)

var life int64

var printTemplate = template.Must(template.New("kubeconfig").Parse(`you can run worker by:

  $ /.../worker --config-file=/.../worker.yaml

with <worker.yaml> include following:
------------------------------------------------

global:
    leader-token: {{.NodeToken}}@{{.LeaderJoinToken}}
network-manager:
    token: {{.NetworkToken}}@{{.IP}}:{{.Port}}

------------------------------------------------
`))

func NewCreateTokenCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-token",
		Usage:     "create worker startup info",
		UsageText: "likuadm [global options] create-token",
		Action:    createToekn,
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "life",
				Usage:       "how many will token info be valid",
				Destination: &life,
				Value:       10,
			},
		},
	}
}

func createToekn(ctx *cli.Context) error {
	client := NewClient()
	if client == nil {
		return fmt.Errorf("fail to init gRPC client")
	}

	defer client.Close()

	c := context.Background()
	var ip string
	var port int
	var nodeToken string
	var nodeBootStrapToken string
	var networkToken string

	if value, err := client.GRPC().BootStrapNetwork(c, &control.BootStrapNetworkRequest{Life: life}); err != nil {
		return err
	} else {
		ip = value.Ip
		port = int(value.GetPort())
		networkToken = value.Token
	}

	if value, err := client.GRPC().NodeToken(c, &control.NoneValue{}); err != nil {
		return err
	} else {
		nodeToken = value.GetToken()
	}

	if value, err := client.GRPC().CreateToken(c, &control.CreateTokenRequest{Life: life, IsAdmin: false}); err != nil {
		return err
	} else {
		nodeBootStrapToken = value.Token.GetToken()
	}

	data := struct {
		NodeToken       string
		LeaderJoinToken string
		NetworkToken    string
		IP              string
		Port            int
	}{
		NodeToken:       nodeToken,
		LeaderJoinToken: nodeBootStrapToken,
		IP:              ip,
		Port:            port,
		NetworkToken:    networkToken,
	}

	printTemplate.Execute(os.Stdout, &data)

	return nil
}
