package fetch_provider

import (
	"lytemp/pkg/schemaparser"
)

func Router(handler IHandler, router schemaparser.IEchoRouter) {
	r := router.Group("")
	r.GET("", handler.FetchProvider)
}
