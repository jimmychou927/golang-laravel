package controller

import "github.com/gin-gonic/gin"

type Home struct {
	Controller
}

func (this *Home) Display(c *gin.Context) {

	this.View(c, "home.index", map[string]interface{}{})
}
