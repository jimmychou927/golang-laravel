package route

import (
	"golang-laravel/app/http/controller"

	"github.com/gin-gonic/gin"
)

func Setup(e *gin.Engine) {

	// e.Use(middleware.IsAuth)

	e.GET("/", (&controller.Home{}).Display)
	e.GET("/home", (&controller.Home{}).Display)
}
