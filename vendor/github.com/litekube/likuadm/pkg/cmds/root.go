package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/litekube/likuadm/pkg/version"
	"github.com/urfave/cli/v2"
)

type AccessConfig struct {
	ip    string
	port  int
	token string
}

var GlobalConfig AccessConfig

var homeDir string = func() string {
	if home, err := os.UserHomeDir(); err != nil {
		return ""
	} else {
		return home
	}
}()

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "likuadm"
	app.Usage = "likuadm, a commond-line tool to control join to litekube."
	app.Version = version.Version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s\n", app.Name, app.Version)
		fmt.Printf("go version %s\n", runtime.Version())
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "ip",
			Usage:       "leader host ip",
			Destination: &GlobalConfig.ip,
			Value:       "127.0.0.1",
		},
		&cli.IntFlag{
			Name:        "port",
			Usage:       "leader control port",
			Destination: &GlobalConfig.port,
			Value:       6442,
		},
		&cli.StringFlag{
			Name:        "token",
			Usage:       "token to authentication",
			Destination: &GlobalConfig.token,
			FilePath:    filepath.Join(homeDir, ".litekube/token"),
		},
	}

	return app
}
