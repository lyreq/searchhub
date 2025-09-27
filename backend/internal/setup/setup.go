package setup

import (
	"log"
	"lytemp/config"
	"lytemp/internal/middleware"
	"lytemp/internal/router"
	"lytemp/internal/utils/pagination"
	"lytemp/pkg/cron"
	"lytemp/pkg/database"
	"lytemp/pkg/lg"
	"lytemp/pkg/mailer"
	"lytemp/pkg/meilisearch"
	"lytemp/pkg/redis"
	"lytemp/pkg/validator"
	"lytemp/pkg/viper"
	"lytemp/pkg/websocket"
	"os"
	"reflect"

	"github.com/Lexographics/go-openapigen"
	"github.com/Lexographics/go-postmangen"
	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/logarweb"
	"github.com/Lexographics/logar/logfilter"
	"github.com/Lexographics/logar/proxy"
	"github.com/Lexographics/logar/proxy/consolelogger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	wsapi "lytemp/internal/api/websocket"
)

func New() (echo *echo.Echo, shutdownFunc func()) {
	cfg := loadConfig()
	config.Set(cfg)

	postmanGen := setupPostmanGen()
	openapigen := setupOpenAPIGen()

	db := setupDatabase(cfg.Database)
	_validator := validator.New(db.Get(), postmanGen, openapigen)
	redis := setupRedis(cfg.Redis)
	meili := setupMeiliSearch(cfg.MeiliSearch)

	cronService := setupCronService()
	setupLogger(cfg.Logar)

	mailerService := mailer.NewGoMailer(cfg.Mail)
	wsService := setupWebsocket()
	echo = setupEcho()
	router := router.New(&cfg, echo, router.Singleton{
		Database:  db,
		Validator: _validator,
		Redis:     redis,
		Meili:     meili,

		CronService: cronService,
		Mailer:      mailerService,
		WSService:   wsService,
	})

	wsRepo := wsapi.NewRepository(db)
	wsApiService := wsapi.NewService(wsRepo)
	wsApiService.Inject(wsService)

	validator.GenerateDocs = true
	router.Routes()

	validator.GenerateDocs = false

	postmanGen.WriteToFile("lytemp-api.postman_collection.json")
	err := openapigen.WriteToFile("www/openapi.json")
	if err != nil {
		log.Println("OpenAPI dosyası oluşturulurken hata:", err)
	}

	cronService.Start()
	return echo, func() {
		Shutdown(router.Singleton)
	}
}

func loadConfig() config.Config {
	cfg := config.Config{}
	if err := viper.Init(viper.Config{
		Path: ".",
		Name: "config",
		Type: viper.YML,
	}, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

func setupWebsocket() *websocket.Payload {
	return websocket.New(make(map[string]chan websocket.Event))
}

func setupDatabase(cfg config.Database) *database.Client {
	client, err := database.ConnectPostgres(cfg.User, cfg.Pass, cfg.Host, uint16(cfg.Port), cfg.Name, cfg.Debug)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Migrate {
		err := client.Migrate()
		if err != nil {
			log.Fatal(err)
		}
	}

	return client
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.Use(echoMiddleware.RequestIDWithConfig(echoMiddleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
		RequestIDHandler: func(c echo.Context, id string) {
			c.Set("requestid", id)
		},
	}))

	e.Any("/logar/*", echo.WrapHandler(logarweb.ServeHTTP(config.Get().App.URL, "/logar", lg.Get())))
	e.GET("/logar", func(c echo.Context) error {
		return c.Redirect(307, "/logar/")
	})
	e.Static("/", "www")

	e.Use(middleware.RequestLogger())
	e.GET("/api-docs/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URLs = []string{"/openapi.json"}
	},
		echoSwagger.DocExpansion("list"),
	))
	e.GET("/api-docs", func(c echo.Context) error {
		return c.Redirect(307, "/api-docs/")
	})

	return e
}

func setupRedis(cfg config.Redis) redis.Client {
	client := redis.New()
	err := client.Connect(cfg.Host, cfg.Port, cfg.Pass)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func setupMeiliSearch(cfg config.MeiliSearch) meilisearch.Client {
	client := meilisearch.New()
	err := client.Connect(cfg.Host, cfg.Port, cfg.Key)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func setupCronService() *cron.CronService {
	return cron.NewService()
}

func setupLogger(cfg config.Logar) logar.App {
	if _, err := os.Stat("logs/logs.db"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	app, err := logar.New(
		logar.WithAppName("lytemp"),
		logar.WithDatabase("logs/logs.db"),

		logar.If(cfg.AdminUsername != "" && cfg.AdminPassword != "",
			logar.WithAdminCredentials(cfg.AdminUsername, cfg.AdminPassword),
		),

		logar.AddModel("Kullanıcı", lg.UserLogs, "fa fa-users"),
		logar.AddModel("Sistem", lg.SystemLogs, "fa fa-server"),
		logar.AddModel("Admin", lg.AdminLogs, "fa-solid fa-user-shield"),
		logar.AddModel("Tüm Loglar", "__all__", "fa-solid fa-file-lines"),

		logar.AddProxy(proxy.NewProxy(consolelogger.New(), logfilter.NewFilter())),
	)

	if err != nil {
		log.Fatal(err)
	}

	lg.Set(app)
	return app
}

func Shutdown(singleton router.Singleton) {
	singleton.CronService.Stop()
}

func setupPostmanGen() *postmangen.PostmanGen {
	postmanGen := postmangen.NewPostmanGen("SEARCHHUB API", "SearchHub API Documentation").
		AddVariable("base_url", "http://localhost:8080").
		AddVariable("token", "YOUR TOKEN HERE")
	return postmanGen
}

func setupOpenAPIGen() *openapigen.OpenAPIGen {
	openapigen := openapigen.New("SEARCHHUB API", "SEARCHHUB API Documentation")
	err := openapigen.RegisterType(reflect.TypeOf(pagination.Pagination{}))
	if err != nil {
		log.Fatal(err)
	}
	openapigen.
		AddGroup("Contents", "^/api/v1/contents(/.*)?$").
		AddGroup("Provider", "^/api/v1/fetch-provider(/.*)?$")

	return openapigen
}
