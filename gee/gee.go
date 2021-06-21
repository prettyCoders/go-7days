package gee

import (
	"log"
	"net/http"
	"strings"
)

//HandlerFunc 定义 gee 处理请求的方法
type HandlerFunc func(*Context)

//Engine
//实现了http包的ServeHTTP方法
//新增路由
//监听请求
//处理请求
type (
	Engine struct {
		*RouterGroup                //将Engine作为最顶层的分组，也就是说Engine拥有RouterGroup所有的能力
		groups       []*RouterGroup //存储所有分组
		router       *router
	}
	RouterGroup struct {
		prefix      string        //前缀
		middlewares []HandlerFunc //支持中间件
		parent      *RouterGroup  //支持嵌套
		engine      *Engine       //所有分组共享一个engine，用于间接地访问engine到各种接口
	}
)

//New Engine构造器
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

//addRouter 添加路由
func (engine *Engine) addRouter(method string, routerPath string, handler HandlerFunc) {
	engine.router.add(method, routerPath, handler)
}

//GET 添加GET路由
func (engine *Engine) GET(routerPath string, handler HandlerFunc) {
	engine.router.add("GET", routerPath, handler)
}

//POST 添加POST路由
func (engine *Engine) POST(routerPath string, handler HandlerFunc) {
	engine.router.add("POST", routerPath, handler)
}

//ServeHTTP http包的ServeHTTP方法实现
func (engine Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	engine.router.handle(c)
}

//Run 启动http服务
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

//Group 创建一个新到路由组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//addRouter 添加路由
func (group *RouterGroup) addRouter(method string, subPattern string, handler HandlerFunc) {
	pattern := group.prefix + subPattern
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.add(method, pattern, handler)
}

//GET 添加GET路由
func (group *RouterGroup) GET(subPattern string, handler HandlerFunc) {
	group.addRouter("GET", subPattern, handler)
}

//POST 添加POST路由
func (group *RouterGroup) POST(subPattern string, handler HandlerFunc) {
	group.addRouter("POST", subPattern, handler)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
