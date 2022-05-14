package cmds

import (
	"context"
	"fmt"

	"github.com/litekube/LiteKube/pkg/leader/runtime/control"
	"github.com/urfave/cli/v2"
)

func NewHealthCommand() *cli.Command {
	return &cli.Command{
		Name:      "health",
		Usage:     "check if server health",
		UsageText: "likuadm [global options] health",
		Action:    health,
	}
}

func health(ctx *cli.Context) error {
	client := NewClient()
	if client == nil {
		return fmt.Errorf("fail to init gRPC client")
	}

	defer client.Close()

	if value, err := client.GRPC().CheckHealth(context.TODO(), &control.NoneValue{}); err != nil {
		fmt.Println(".....")
		return err
	} else {
		fmt.Println(value.Message)
		return nil
	}
}
