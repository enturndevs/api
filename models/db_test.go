package models

import (
	l "github.com/enturndevs/enturn-service/core/logger"
)

func init() {
	l.Init(true)
	InitDB()
	InitRedis()

	x.Exec("TRUNCATE TABLE user")
	r.FlushAll()
}
