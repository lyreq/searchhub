package websocket

import (
	"lytemp/pkg/schemaparser"
)

func Router(wsHandler *Handler, router schemaparser.IEchoRouter) {
	router.GET("/", wsHandler.HandleWebSocket)
}
