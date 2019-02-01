package cmd

import (
	"github.com/stanlyliao/logger"
	"github.com/urfave/cli"
)

// CmdWeb command web server.
var CmdWeb = cli.Command{
	Name:   "web",
	Usage:  "Start enturn web server",
	Action: runWeb,
	Flags:  []cli.Flag{},
}

func runWeb(ctx *cli.Context) error {
	logger.Info("Enturn web service v" + ctx.App.Version)
	return nil
}
