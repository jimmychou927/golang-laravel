package database

import (
	"database/sql"

	"extension/database/mysql"
	"golang-laravel/config"
)

type Connection interface {
	Query(query string, args ...interface{}) ([]map[string]interface{}, *sql.Rows)
	Exec(query string, args ...interface{}) sql.Result
	InitDB(cfg map[string]config.Database)
}

func GetConnectionByDriver(driver string) Connection {
	switch driver {
	case "mysql":
		return mysql.GetMysqlDB()
	default:
		panic("driver not found!")
	}
}

func GetConnection() Connection {
	return GetConnectionByDriver(config.Get().DATABASE[0].DRIVER)
}

func Query(query string, args ...interface{}) ([]map[string]interface{}, *sql.Rows) {
	return GetConnection().Query(query, args...)
}

func Exec(query string, args ...interface{}) sql.Result {
	return GetConnection().Exec(query, args...)
}
