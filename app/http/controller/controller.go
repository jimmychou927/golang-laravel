package controller

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

var ControllerMap map[string]reflect.Type

type Controller struct {
}

func (this *Controller) Redirect(c *gin.Context, path string) {

	c.Redirect(http.StatusMovedPermanently, path)
}

func (this *Controller) View(c *gin.Context, name string, assign map[string]interface{}) {

	c.HTML(http.StatusOK, name, assign)
}

func register(ctrl interface{}) {

	if ControllerMap == nil {
		ControllerMap = make(map[string]reflect.Type)
	}

	t := reflect.TypeOf(ctrl)
	ControllerMap[t.Name()] = t
}
