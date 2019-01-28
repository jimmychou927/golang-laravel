package config

type Cookie struct {
	PATH     string
	LIFTTIME int
	SECURE   bool
	HTTPONLY bool
}

type Database struct {
	HOST         string
	PORT         string
	USER         string
	PWD          string
	NAME         string
	MAX_IDLE_CON int
	MAX_OPEN_CON int
	DRIVER       string
	FILE         string
}

type Server struct {
	PORT int
}

type View struct {
	PATH string
}

type Config struct {
	COOKIE   Cookie
	DATABASE []Database
	SERVER   Server
	VIEW     View
}

var (
	globalCfg Config
)

func init() {

	globalCfg = Config{
		COOKIE: defaultCookieCfg,
		DATABASE: []Database{
			defaultDatabaseCfg,
		},
		SERVER: defaultServerCfg,
		VIEW:   defaultViewCfg,
	}
}

func Get() Config {

	return globalCfg
}
