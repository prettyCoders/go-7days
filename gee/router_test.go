package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.add("GET", "/", nil)
	r.add("GET", "/hello/:name", nil)
	r.add("GET", "/hello/b/c", nil)
	r.add("GET", "/hi/:name", nil)
	r.add("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parseRouterPath("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parseRouterPath("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parseRouterPath("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parseRouterPath failed")
	}
}

func TestGetRouter(t *testing.T) {
	r := newTestRouter()
	n, params := r.get("GET", "/hello/geektutu")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.routerPath != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if params["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.routerPath, params["name"])
}
