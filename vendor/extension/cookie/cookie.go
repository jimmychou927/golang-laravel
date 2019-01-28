package cookie

import (
	"golang-laravel/config"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context, name string) string {

	cookie, err := c.Cookie(name)

	if err != nil {
		return ""
	}

	return cookie
}

func Put(c *gin.Context, name string, value string) {

	c.SetCookie(name, value, config.Get().COOKIE.LIFTTIME*60, config.Get().COOKIE.PATH, location.Get(c).Host, config.Get().COOKIE.SECURE, config.Get().COOKIE.HTTPONLY)
}
