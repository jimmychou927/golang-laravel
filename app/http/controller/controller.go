package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func (this *Controller) Redirect(c *gin.Context, path string) {

	c.Redirect(http.StatusMovedPermanently, path)
}

func (this *Controller) View(c *gin.Context, name string, assign map[string]interface{}) {

	c.HTML(http.StatusOK, name, assign)
}
