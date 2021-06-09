package gee

import (
	"log"
	"net/http"
)

//router 路由
type router struct {
	handlers map[string]HandlerFunc
}

//newRouter 创建路由对象
func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

//add 添加路由
func (r *router) add(method string, pattern string, handler HandlerFunc) {
	log.Printf("Router:%4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

//handle 处理路由
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
	}
}
