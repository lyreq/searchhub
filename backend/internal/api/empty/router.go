package empty

import (
	"lytemp/internal/middleware"
	"lytemp/pkg/schemaparser"
)

func Router(handler *Handler, router schemaparser.IEchoRouter) {
	protected := router.Group("", middleware.Authentication)
	protected.POST("/example", handler.ExampleRoute)
}
