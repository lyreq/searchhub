package contents

import (
	"context"
	"lytemp/internal/domain/dtos/responses"
	"lytemp/internal/domain/models"
	"time"

	"github.com/labstack/echo/v4"
)

type IRepository interface {
	Search(ctx context.Context, Query string, Type string, Sort string, Page int, PerPage int) ([]models.Content, int64, error)
	GetByID(ctx context.Context, id uint) (models.Content, error)
}

type IService interface {
	WithCache(get func(ctx context.Context, key string) (string, error), set func(ctx context.Context, key string, value string, ttl time.Duration) error, ttl time.Duration) *Service
	List(ctx context.Context, Query string, Type string, Sort string, Page int, PerPage int) (responses.ListContentEnvelopeResponse, error)
	Detail(ctx context.Context, id uint, includeRaw bool) (responses.ShowContentResponse, error)
	makeCacheKey(Query string, Type string, Sort string, Page int, PerPage int) string
}

type IHandler interface {
	Index(c echo.Context) error
	Show(c echo.Context) error
}
