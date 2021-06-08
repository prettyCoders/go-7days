package gee

import (
	"fmt"
	"net/http"
)

//HandlerFunc 定义 gee 处理请求的方法
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

//Engine 实现了http包的ServeHTTP方法
type Engine struct {
	router map[string]HandlerFunc
}

//New Engine构造器
func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

//addRouter 添加路由
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

//Get 添加Get路由
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

//Post 添加Post路由
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("Post", pattern, handler)
}

//Run 启动http服务
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, r)
	} else {
		_, _ = fmt.Fprintf(w, "404 NOT FOUND:%s\n", r.URL)
	}
}
