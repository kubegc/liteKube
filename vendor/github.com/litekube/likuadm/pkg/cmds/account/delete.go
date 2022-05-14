package account

import (
	"context"
	"fmt"

	"github.com/litekube/LiteKube/pkg/leader/runtime/control"
	"github.com/litekube/likuadm/pkg/cmds"
	"github.com/urfave/cli/v2"
)

var token string

func NewDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete-account",
		Usage:     "delete leader control service account",
		UsageText: "likuadm [global options] delete-account [options]",
		Action:    delete,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "token",
				Usage:       "token to mark account to be delete",
				Destination: &token,
				Required:    true,
			},
		},
	}
}

func delete(ctx *cli.Context) error {
	client := cmds.NewClient()
	if client == nil {
		return fmt.Errorf("fail to init gRPC client")
	}

	defer client.Close()

	if value, err := client.GRPC().DeleteToken(context.Background(), &control.TokenString{Token: token}); err != nil {
		return err
	} else {
		fmt.Printf("status code: %d\nmessage: %s\n\n", value.GetStatusCode(), value.GetMessage())
	}

	return nil
}
