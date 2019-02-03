package main

import (
	"os"

	"github.com/enturndevs/enturn-service/models"

	"github.com/enturndevs/enturn-service/cmd"
	c "github.com/enturndevs/enturn-service/core/config"
	l "github.com/enturndevs/enturn-service/core/logger"

	"github.com/stanlyliao/logger"
	"github.com/urfave/cli"
)

// Version enturn service ver.
var Version = "0.0.1"

func main() {
	c.Init()
	l.Init()
	models.InitDB()
	models.InitRedis()

	app := cli.NewApp()
	app.Name = "Enturn"
	app.Usage = "Micro Service"
	app.Version = Version
	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}
	app.Flags = append(app.Flags, cmd.CmdWeb.Flags...)
	app.Action = cmd.CmdWeb.Action

	if err := app.Run(os.Args); err != nil {
		logger.Panic(err)
	}
}
