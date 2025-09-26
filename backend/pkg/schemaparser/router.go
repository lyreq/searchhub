package schemaparser

import (
	"fmt"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

type IEchoRouter interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Use(m ...echo.MiddlewareFunc)
	Group(path string, m ...echo.MiddlewareFunc) IEchoRouter
}

type DocGenRouter struct {
	echo  *echo.Echo
	group *echo.Group
}

func FromGroup(echo *echo.Echo, group *echo.Group) IEchoRouter {
	return &DocGenRouter{echo: echo, group: group}
}

func (r *DocGenRouter) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	route := r.group.GET(path, h, m...)
	r.Generate(h, route)
	return route
}

func (r *DocGenRouter) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	route := r.group.POST(path, h, m...)
	r.Generate(h, route)
	return route
}

func (r *DocGenRouter) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	route := r.group.PUT(path, h, m...)
	r.Generate(h, route)
	return route
}

func (r *DocGenRouter) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	route := r.group.DELETE(path, h, m...)
	r.Generate(h, route)
	return route
}

func (r *DocGenRouter) Use(m ...echo.MiddlewareFunc) {
	r.group.Use(m...)
}

func (r *DocGenRouter) Group(path string, m ...echo.MiddlewareFunc) IEchoRouter {
	return FromGroup(r.echo, r.group.Group(path, m...))
}

func (r *DocGenRouter) Generate(h echo.HandlerFunc, route *echo.Route) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Failed to generate docs for %s %s: %s\n", route.Method, route.Path, r)
		}
	}()

	c := r.echo.AcquireContext()
	defer r.echo.ReleaseContext(c)

	c.SetRequest(httptest.NewRequest(route.Method, route.Path, nil))
	c.SetResponse(echo.NewResponse(httptest.NewRecorder(), r.echo))

	h(c)
}
