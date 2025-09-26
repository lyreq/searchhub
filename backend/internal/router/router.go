package router

import (
	"lytemp/config"
	"lytemp/internal/api/empty"
	ws "lytemp/internal/api/websocket"
	"lytemp/pkg/cron"
	"lytemp/pkg/database"
	"lytemp/pkg/iyzico"
	"lytemp/pkg/mailer"
	"lytemp/pkg/meilisearch"
	notificationService "lytemp/pkg/notification"
	"lytemp/pkg/onesignal"
	"lytemp/pkg/redis"
	"lytemp/pkg/schemaparser"
	"lytemp/pkg/telegram"
	"lytemp/pkg/validator"
	"lytemp/pkg/websocket"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
)

type Router struct {
	Config    *config.Config
	Echo      *echo.Echo
	Singleton Singleton
}

type Singleton struct {
	Database     *database.Client
	Validator    validator.Validator
	Redis        redis.Client
	Meili        meilisearch.Client
	CronService  *cron.CronService
	GCSClient    *storage.Client
	GCSBucket    string
	Mailer       mailer.Mailer
	Telegram     *telegram.Telegram
	WSService    *websocket.Payload
	Iyzico       iyzico.Iyzico
	OneSignal    *onesignal.OneSignalService
	Notification *notificationService.Service
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

	emptyRepository := empty.NewRepository(r.Singleton.Database)
	emptyService := empty.NewService(emptyRepository)
	emptyHandler := empty.NewHandler(emptyService, r.Singleton.Validator)

	wsRepository := ws.NewRepository(r.Singleton.Database)
	wsService := ws.NewService(wsRepository)
	wsHandler := ws.NewHandler(wsService, r.Singleton.Validator)

	//
	// ****************************************** //
	// **************** Routers ***************** //
	// ****************************************** //

	wsService.Inject(r.Singleton.WSService)
	ws.Router(wsHandler, apiv1.Group("/ws"))

	empty.Router(emptyHandler, apiv1.Group("/empty"))

}
