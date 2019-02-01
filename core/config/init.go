package config

import (
	"github.com/spf13/viper"
)

// Init 初始化
func Init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
