package gee

import (
	"net/http"
	"strings"
)

//router 路由
type router struct {
	trieRoots map[string]*node       //存储每种请求方式(GET、POST)的Trie 树根节点
	handlers  map[string]HandlerFunc //存储每种请求方式的 HandlerFunc
}

//newRouter 创建路由对象
func newRouter() *router {
	return &router{
		trieRoots: make(map[string]*node),
		handlers:  make(map[string]HandlerFunc),
	}
}

//parseRouterPath 解析路由路径,根据/分割路径，返回路由组成部分的slice
func parseRouterPath(routerPath string) []string {
	pathValues := strings.Split(routerPath, "/")
	parts := make([]string, 0)
	for _, value := range pathValues {
		if value != "" {
			parts = append(parts, value)
			if value[0] == '*' {
				break
			}
		}
	}
	return parts
}

//add 添加路由
func (r *router) add(method string, routerPath string, handler HandlerFunc) {
	routerValues := parseRouterPath(routerPath)
	key := method + "-" + routerPath
	_, ok := r.trieRoots[method]
	if !ok {
		r.trieRoots[method] = &node{}
	}
	r.trieRoots[method].insert(routerPath, routerValues, 0)
	r.handlers[key] = handler
}

//get 匹配node，同时返回请求参数
func (r *router) get(method string, path string) (*node, map[string]string) {
	searchParts := parseRouterPath(path)
	params := make(map[string]string)
	root, ok := r.trieRoots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parseRouterPath(n.routerPath)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

//handle 处理路由
func (r *router) handle(c *Context) {
	n, params := r.get(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.routerPath
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
		})
	}
	c.Next()
}
