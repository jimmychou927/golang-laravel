package config

var (
	defaultDatabaseCfg = Database{
		HOST:         "localhost",
		PORT:         "3306",
		USER:         "golang-laravel",
		PWD:          "golang-laravel",
		NAME:         "golang-laravel",
		MAX_IDLE_CON: 50,
		MAX_OPEN_CON: 150,
		DRIVER:       "mysql",
		FILE:         "",
	}
)
