package fetch_provider

import (
	"lytemp/internal/domain/dtos/requests"
	"lytemp/internal/domain/dtos/responses"
	"lytemp/pkg/validator"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	validator validator.Validator
	service   IService
}

func NewHandler(service IService, validator validator.Validator) *Handler {
	return &Handler{service: service, validator: validator}
}

func (h *Handler) FetchProvider(c echo.Context) error {
	ctx := c.Request().Context()

	var q requests.FetchProviderRequest

	if err := h.validator.ParseAndValidate(c, &q); err != nil {
		return c.JSON(400, responses.Error(err, err.Error()))
	}

	// provider da mock olduğu için varsayılanlar
	q.Pages = 1
	q.PerPage = 10

	out, err := h.service.FetchAndStore(ctx, q.Query, q.Pages, q.PerPage)
	if err != nil {
		return c.JSON(500, responses.Error(err, err.Error()))
	}
	return c.JSON(200, responses.Success(out, "providers fetched & stored"))
}
