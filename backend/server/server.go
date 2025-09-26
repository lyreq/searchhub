package server

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

type Server interface {
	Start()
	Stop()
	RegisterValidations()
}

type server struct {
	echo      *echo.Echo
	port      string
	validator *validator.Validate
	db        *gorm.DB
}

type Handler interface {
	AddRoutes(e *echo.Echo)
}

// @title			lytemp BACKEND API
// @version		1.0
// @description	Backend service for lytemp project
func NewServer(port string, handlers []Handler, db *gorm.DB) Server {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	time.Local, _ = time.LoadLocation("Europe/Istanbul")
	v := validator.New()
	e.Validator = &Validator{Validator: v}
	for _, h := range handlers {
		h.AddRoutes(e)
	}
	return &server{
		port:      port,
		echo:      e,
		validator: v,
		db:        db,
	}
}

func (s *server) Start() {
	s.echo.Logger.Fatal(s.echo.Start(":" + s.port))
}

func (s *server) Stop() {
	var ctx context.Context
	s.echo.Shutdown(ctx)
}

type Validator struct {
	Validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	cv.Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return cv.Validator.Struct(i)
}
