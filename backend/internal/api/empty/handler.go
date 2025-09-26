package empty

import (
	"lytemp/internal/domain/dtos/responses"
	"lytemp/pkg/validator"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	validator validator.Validator
	service   *Service
}

func NewHandler(service *Service, validator validator.Validator) *Handler {
	return &Handler{service: service, validator: validator}
}

func (h *Handler) ExampleRoute(c echo.Context) error {
	ctx := c.Request().Context()

	var rq CreateCity
	if err := h.validator.ParseAndValidate(c, &rq); err != nil {
		return c.JSON(400, responses.Error(err, err.Error()))
	}

	resp, err := h.service.Example(ctx, rq)
	if err != nil {
		return c.JSON(400, responses.Error(err, err.Error()))
	}
	return c.JSON(200, responses.Success(resp, "example servisten gelen veri"))
}
