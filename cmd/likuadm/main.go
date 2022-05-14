package main

import (
	"fmt"
	"os"

	"github.com/litekube/likuadm/pkg/cmds"
	"github.com/litekube/likuadm/pkg/cmds/account"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cmds.NewApp()
	app.Commands = []*cli.Command{
		cmds.NewHealthCommand(),
		cmds.NewCreateTokenCommand(),
		account.NewCreateCommand(),
		account.NewListCommand(),
		account.NewDeleteCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error options: %s\n", err.Error())
		os.Exit(-1)
	}
}
