package gee

import (
	"net/http"
)

//HandlerFunc 定义 gee 处理请求的方法
type HandlerFunc func(*Context)

//Engine 实现了http包的ServeHTTP方法
type Engine struct {
	router *router
}

//New Engine构造器
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

//addRouter 添加路由
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router.add(method, pattern, handler)
}

//GET 添加GET路由
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.add("GET", pattern, handler)
}

//POST 添加POST路由
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.router.add("POST", pattern, handler)
}

//ServeHTTP http包的ServeHTTP方法实现
func (engine Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.router.handle(c)
}

//Run 启动http服务
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
