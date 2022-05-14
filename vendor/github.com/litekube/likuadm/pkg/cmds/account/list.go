package account

import (
	"context"
	"fmt"

	"github.com/litekube/LiteKube/pkg/leader/runtime/control"
	"github.com/litekube/likuadm/pkg/cmds"
	"github.com/urfave/cli/v2"
)

func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:      "list-accounts",
		Usage:     "list all leader control service accounts",
		UsageText: "likuadm [global options] list-accounts",
		Action:    list,
	}
}

func list(ctx *cli.Context) error {
	client := cmds.NewClient()
	if client == nil {
		return fmt.Errorf("fail to init gRPC client")
	}

	defer client.Close()

	if value, err := client.GRPC().QueryTokens(context.Background(), &control.NoneValue{}); err != nil {
		return err
	} else {
		fmt.Printf("status code: %d\nmessage: %s\n\n", value.GetStatusCode(), value.GetMessage())
		if value.GetStatusCode() >= 200 && value.GetStatusCode() < 300 {
			for _, v := range value.GetTokenList() {
				if v.GetLife() < 0 {
					fmt.Printf("token: %s\ncreate time: %s\nlife: permanent\nadmin: %t\nvalid now: always\n\n", v.GetToken(), v.GetCreateTime(), v.GetIsAdmin())
				} else {
					fmt.Printf("token: %s\ncreate time: %s\nlife: %d minutes\nadmin: %t\nvalid now: %t\n\n", v.GetToken(), v.GetCreateTime(), v.GetLife(), v.GetIsAdmin(), v.GetValid())
				}
			}
		}
	}

	return nil
}
