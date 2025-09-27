package fetch_provider

import (
	"context"
	"time"

	"lytemp/internal/domain/dtos/responses"
	"lytemp/internal/domain/models"

	"github.com/labstack/echo/v4"
)

type IRepository interface {
	UpsertContent(ctx context.Context, c *models.Content) (bool, error)
	UpsertProviderSync(ctx context.Context, name string, format models.ProviderFormat, rateLimit int, status models.ProviderStatus, lastSync time.Time, lastErr *string) error
}

type IService interface {
	FetchAndStore(ctx context.Context, query string, pages, perPage int) (responses.FetchResult, error)
}

type IHandler interface {
	FetchProvider(c echo.Context) error
}
