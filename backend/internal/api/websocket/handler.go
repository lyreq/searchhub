package websocket

import (
	"lytemp/pkg/validator"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	validator validator.Validator
	service   *Service
}

func NewHandler(service *Service, validator validator.Validator) *Handler {
	return &Handler{service: service, validator: validator}
}

func (h *Handler) HandleWebSocket(c echo.Context) error {
	clientID := c.QueryParam("client_id")
	if clientID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "client_id gerekli"})
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Client'ı kaydet
	eventChan := h.service.RegisterClient(clientID)
	defer h.service.UnregisterClient(clientID)

	// Event'leri dinle ve WebSocket üzerinden gönder
	for event := range eventChan {
		if err := ws.WriteJSON(event); err != nil {
			return err
		}
	}

	return nil
}
