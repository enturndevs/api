package models

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/spf13/viper"
	"github.com/stanlyliao/logger"

	// load mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	x *xorm.Engine
)

// InitDB db initial
func InitDB() {
	var err error

	x, err = xorm.NewEngine("mysql", getDSN())
	if err != nil {
		logger.Panic(err)
	}
	x.SetMaxIdleConns(0)
}

// // getDSN get data source name
func getDSN() string {
	username := viper.GetString("db.username")
	if username == "" {
		username = "root"
	}
	password := viper.GetString("db.password")
	if password == "" {
		password = "root"
	}
	host := viper.GetString("db.host")
	if host == "" {
		host = "127.0.0.1"
	}
	port := viper.GetString("db.port")
	if port == "" {
		port = "3306"
	}
	dbname := viper.GetString("db.name")
	if dbname == "" {
		dbname = "enture"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=%s", username, password, host, port, dbname, "Asia%2FTaipei")
}
