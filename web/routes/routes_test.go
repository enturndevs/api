package routes

import (
	l "github.com/enturndevs/enturn-service/core/logger"
	"github.com/enturndevs/enturn-service/models"
)

func init() {
	l.Init(true)
	models.InitDB()
	models.InitRedis()

	initServer()
	initRoutes()
}
