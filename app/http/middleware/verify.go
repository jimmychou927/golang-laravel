package middleware

import (
	"extension/cookie"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAuth(c *gin.Context) {

	// before request statement
	if cookie.Get(c, "token") != "..." {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}

	c.Next()

	// after request statement
}
