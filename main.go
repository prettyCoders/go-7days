package main

import (
	"github.com/prettyCoders/go-7days/gee"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>hello gee</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s ,you are at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
