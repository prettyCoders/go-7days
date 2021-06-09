package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//H 别名，构建JSON数据时，显得更简洁
type H map[string]interface{}

//Context 随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由 Context 承载
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

//newContext 创建新的Context
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

//PostForm 从POST表单中取数据
func (c *Context) PostForm(key string) string {
	return c.Req.PostFormValue(key)
}

//Query 从GET请求路径中取数据
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//Status 设置响应状态码
func (c *Context) Status(statusCode int) {
	c.StatusCode = statusCode
	c.Writer.WriteHeader(statusCode)
}

//SetHeader 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//String 响应文本
func (c *Context) String(statusCode int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(statusCode)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//JSON 响应JSON
func (c *Context) JSON(statusCode int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(statusCode)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

//Data 直接响应字节数组数据
func (c *Context) Data(statusCode int, data []byte) {
	c.Status(statusCode)
	c.Writer.Write(data)
}

//HTML 响应HTML
func (c *Context) HTML(statusCode int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(statusCode)
	c.Writer.Write([]byte(html))
}
