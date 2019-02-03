package models

import (
	"database/sql"
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/spf13/viper"
	"github.com/stanlyliao/logger"

	// load mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Engine represents a xorm engine or session.
type Engine interface {
	Table(tableNameOrBean interface{}) *xorm.Session
	Count(...interface{}) (int64, error)
	Decr(column string, arg ...interface{}) *xorm.Session
	Delete(interface{}) (int64, error)
	Exec(...interface{}) (sql.Result, error)
	Find(interface{}, ...interface{}) error
	Get(interface{}) (bool, error)
	ID(interface{}) *xorm.Session
	In(string, ...interface{}) *xorm.Session
	Incr(column string, arg ...interface{}) *xorm.Session
	Insert(...interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	Iterate(interface{}, xorm.IterFunc) error
	Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *xorm.Session
	SQL(interface{}, ...interface{}) *xorm.Session
	Where(interface{}, ...interface{}) *xorm.Session
}

var (
	x      *xorm.Engine
	tables []interface{}
)

// InitDB initial db and sync schema
func InitDB() {
	var err error

	x, err = xorm.NewEngine("mysql", getDSN())
	if err != nil {
		logger.Panic(err)
	}
	x.SetMaxIdleConns(0)

	tables = []interface{}{
		new(User),
	}

	if err = x.Sync2(tables...); err != nil {
		logger.Panic(err)
	}
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
		dbname = "enturn"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=%s", username, password, host, port, dbname, "Asia%2FTaipei")
}
