package main

import (
	c "github.com/enturndevs/enturn-service/core/config"
	l "github.com/enturndevs/enturn-service/core/logger"
	"github.com/stanlyliao/logger"
)

// Version enturn service ver.
var Version = "0.0.1"

func main() {
	c.Init()
	l.Init()

	logger.Info("Enturn service v" + Version)
}
