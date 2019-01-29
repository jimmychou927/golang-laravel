package reflection

import (
	"fmt"
	"golang-laravel/app/http/controller"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebRoute interface {
	ALL(string, string) WebRoute
	GET(string, string) WebRoute
	POST(string, string) WebRoute
	PUT(string, string) WebRoute
	DELETE(string, string) WebRoute
	PATCH(string, string) WebRoute
	HEAD(string, string) WebRoute
	OPTIONS(string, string) WebRoute
}

type Handler struct {
	controller reflect.Type
	method     string
}

func (h *Handler) handle(request *gin.Context) {
	ctrlValue := reflect.New(h.controller)
	m := ctrlValue.MethodByName(h.method)
	if !m.IsValid() {
		fmt.Printf("[ERROR] handler: find no method: %+v\n", h.controller.String()+"."+h.method)
		return
	}
	m.Call([]reflect.Value{reflect.ValueOf(request)})
}

type Router struct {
	Engine *gin.Engine
}

func (r *Router) binding(method string, url string, handler string) WebRoute {
	handlerData := strings.Split(handler, "@")

	ctrlType := controller.ControllerMap[handlerData[0]]
	if ctrlType == nil {
		fmt.Printf("[ERROR] binding fail: unregist controller: %+v\n", handlerData[0])
		return r
	}

	h := Handler{controller: ctrlType, method: handlerData[1]}

	args := make([]reflect.Value, 2)
	args[0] = reflect.ValueOf(url)
	args[1] = reflect.ValueOf(h.handle)

	reflect.ValueOf(r.Engine).MethodByName(method).Call(args)
	return r
}

func (r *Router) ALL(url string, handler string) WebRoute {
	return r.binding("ALL", url, handler)
}

func (r *Router) GET(url string, handler string) WebRoute {
	return r.binding("GET", url, handler)
}

func (r *Router) POST(url string, handler string) WebRoute {
	return r.binding("POST", url, handler)
}

func (r *Router) PUT(url string, handler string) WebRoute {
	return r.binding("PUT", url, handler)
}

func (r *Router) DELETE(url string, handler string) WebRoute {
	return r.binding("DELETE", url, handler)
}

func (r *Router) PATCH(url string, handler string) WebRoute {
	return r.binding("PATCH", url, handler)
}

func (r *Router) HEAD(url string, handler string) WebRoute {
	return r.binding("HEAD", url, handler)
}

func (r *Router) OPTIONS(url string, handler string) WebRoute {
	return r.binding("OPTIONS", url, handler)
}
