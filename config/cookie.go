package config

var (
	defaultCookieCfg = Cookie{
		PATH:     "/",
		LIFTTIME: 60,
		SECURE:   false,
		HTTPONLY: true,
	}
)
