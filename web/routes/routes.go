package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	r *gin.Engine
)

// Response represents base response
type Response struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error_msg,omitempty"`
}

// Request represents base request
type Request map[string]interface{}

func initServer() {
	if viper.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r = gin.Default()
}

func initRoutes() {
	initUserRoutes()
}

// Start start web server
func Start() {
	initRoutes()
	r.Run(viper.GetString("server.port"))
}

func responseError(c *gin.Context, status int, msg string) {
	c.JSON(status, Response{
		Success:  false,
		ErrorMsg: msg,
	})
}
