package routedumper

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
)

type Route struct {
	Method string
	Path   string
}

func LoadRoutes(e *echo.Echo) []Route {
	routes := []Route{}

	for _, route := range e.Routes() {
		if !slices.Contains([]string{"GET", "POST", "PUT", "DELETE"}, route.Method) {
			continue
		}

		routes = append(routes, Route{
			Method: route.Method,
			Path:   route.Path,
		})
	}

	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			order := map[string]int{"POST": 0, "GET": 1, "PUT": 2, "DELETE": 3}
			return order[routes[i].Method] < order[routes[j].Method]
		}
		return routes[i].Path < routes[j].Path
	})

	return routes
}

func Dump(e *echo.Echo, color bool) {
	routes := LoadRoutes(e)

	getColor := func(method string) string {
		if !color {
			return ""
		}

		switch method {
		case "GET":
			return "\033[32m"
		case "POST":
			return "\033[33m"
		case "PUT":
			return "\033[34m"
		case "DELETE":
			return "\033[31m"
		}
		return ""
	}
	getColorReset := func() string {
		if !color {
			return ""
		}
		return "\033[0m"
	}

	for _, route := range routes {
		fmt.Printf("route: %s%s%s %s%s\n", getColor(route.Method), route.Method, strings.Repeat(" ", 8-len(route.Method)), route.Path, getColorReset())
	}
}
