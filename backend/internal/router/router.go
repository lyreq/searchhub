package router

import (
	"context"
	"lytemp/config"
	"lytemp/internal/api/contents"
	"lytemp/internal/api/fetch_provider"
	"lytemp/pkg/cron"
	"lytemp/pkg/database"
	"lytemp/pkg/redis"
	"lytemp/pkg/schemaparser"
	"lytemp/pkg/validator"
	"time"

	"github.com/labstack/echo/v4"
)

type Router struct {
	Config    *config.Config
	Echo      *echo.Echo
	Singleton Singleton
}

type Singleton struct {
	Database    *database.Client
	Validator   validator.Validator
	Redis       redis.Client
	CronService *cron.CronService
}

func New(config *config.Config, echo *echo.Echo, singleton Singleton) *Router {
	return &Router{
		Config:    config,
		Echo:      echo,
		Singleton: singleton,
	}
}

func (r *Router) Routes() {
	apiv1 := schemaparser.FromGroup(r.Echo, r.Echo.Group("/api/v1"))
	// apiv2 := r.Echo.Group("/api/v2")

	r.Echo.Static("/assets", "assets")
	r.Echo.Static("/uploads", "uploads")

	//
	// ****************************************** //
	// ************** Controllers *************** //
	// ****************************************** //

	// emptyRepository := empty.NewRepository(r.Singleton.Database)
	// emptyService := empty.NewService(emptyRepository)
	// emptyHandler := empty.NewHandler(emptyService, r.Singleton.Validator)

	fetchProviderRepository := fetch_provider.NewRepository(r.Singleton.Database)
	fetchProviderService := fetch_provider.NewService(fetchProviderRepository)
	fetchProviderHandler := fetch_provider.NewHandler(fetchProviderService, r.Singleton.Validator)

	contentsRepo := contents.NewRepository(r.Singleton.Database)
	contentsSvc := contents.NewService(contentsRepo)
	contentsSvc.WithCache(
		func(ctx context.Context, key string) (string, error) {
			return r.Singleton.Redis.Get(key)
		},
		func(ctx context.Context, key string, value string, ttl time.Duration) error {
			return r.Singleton.Redis.Set(key, value, ttl)
		},
		60*time.Second,
	)

	contentsHdl := contents.NewHandler(contentsSvc, r.Singleton.Validator)

	// wsRepository := ws.NewRepository(r.Singleton.Database)
	// wsService := ws.NewService(wsRepository)
	// wsHandler := ws.NewHandler(wsService, r.Singleton.Validator)

	//
	// ****************************************** //
	// **************** Routers ***************** //
	// ****************************************** //

	// wsService.Inject(r.Singleton.WSService)
	// ws.Router(wsHandler, apiv1.Group("/ws"))

	// empty.Router(emptyHandler, apiv1.Group("/empty"))

	fetch_provider.Router(fetchProviderHandler, apiv1.Group("/fetch-provider"))
	contents.Router(contentsHdl, apiv1.Group("/contents"))

}
