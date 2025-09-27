package contents

import (
	"lytemp/internal/domain/dtos/requests"
	"lytemp/internal/domain/dtos/responses"
	"lytemp/pkg/validator"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	v validator.Validator
	s IService
}

func NewHandler(s IService, v validator.Validator) *Handler {
	return &Handler{v: v, s: s}
}

func (h *Handler) Index(c echo.Context) error {
	ctx := c.Request().Context()

	var rq requests.ListContentRequest
	if err := h.v.ParseAndValidate(c, &rq); err != nil {
		return c.JSON(400, responses.Error(err, err.Error()))
	}

	resp, err := h.s.List(ctx, rq.Query, rq.Type, rq.Sort, rq.Page, rq.PerPage)
	if err != nil {
		return c.JSON(500, responses.Error(err, err.Error()))
	}
	return c.JSON(200, responses.Success(resp, "Contents listed"))
}

func (h *Handler) Show(c echo.Context) error {
	ctx := c.Request().Context()

	var rq requests.ContentShowRequest

	if err := h.v.ParseAndValidate(c, &rq); err != nil {
		return c.JSON(400, responses.Error(err, "invalid parameters"))
	}

	includeRaw := false
	if rq.IncludeRaw != nil {
		includeRaw = *rq.IncludeRaw
	}

	out, err := h.s.Detail(ctx, rq.ID, includeRaw)
	if err != nil {
		return c.JSON(500, responses.Error(err, err.Error()))
	}
	return c.JSON(200, responses.Success(out, "contents details listed"))
}
