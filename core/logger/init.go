package logger

import (
	"github.com/spf13/viper"
	"github.com/stanlyliao/logger"
)

// Init 初始化
func Init(isTest ...bool) {
	if len(isTest) > 0 {
		return
	}
	if viper.GetBool("debug") {
		logger.SetLevel(logger.DebugLevel)
		logger.SetLog("./logs", "service")
	} else {
		logger.SetLog("./logs", "service")
	}
}
