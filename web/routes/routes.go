package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	r *gin.Engine
)

func initServer() {
	if viper.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r = gin.Default()
}

func initRoutes() {

}

// Start start web server
func Start() {
	initRoutes()
	r.Run(viper.GetString("server.port"))
}
