package account

import (
	"context"
	"fmt"

	"github.com/litekube/LiteKube/pkg/leader/runtime/control"
	"github.com/litekube/likuadm/pkg/cmds"
	"github.com/urfave/cli/v2"
)

var life int64
var isAdmin bool

func NewCreateCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-account",
		Usage:     "create leader control service account",
		UsageText: "likuadm [global options] create-account [options]",
		Action:    create,
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "life",
				Usage:       "how many will token info be valid",
				Destination: &life,
				Value:       10,
			},
			&cli.BoolFlag{
				Name:        "admin",
				Usage:       "is one administrator account",
				Destination: &isAdmin,
				Value:       false,
			},
		},
	}
}

func create(ctx *cli.Context) error {
	client := cmds.NewClient()
	if client == nil {
		return fmt.Errorf("fail to init gRPC client")
	}

	defer client.Close()

	if value, err := client.GRPC().CreateToken(context.Background(), &control.CreateTokenRequest{Life: life, IsAdmin: isAdmin}); err != nil {
		return err
	} else {
		fmt.Printf("status code: %d\nmessage: %s\n\n", value.GetStatusCode(), value.GetMessage())
		if value.GetStatusCode() >= 200 && value.GetStatusCode() < 300 {
			fmt.Printf("token: %s\ncreate time: %s\nlife: %d minutes\n\n", value.GetToken().GetToken(), value.GetToken().GetCreateTime(), value.GetToken().GetLife())
		}
	}

	return nil
}
