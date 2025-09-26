package main

import (
	"context"
	"fmt"
	"log"
	"lytemp/config"
	"lytemp/internal/setup"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	echo, shutdownFunc := setup.New()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := echo.Start(fmt.Sprintf(":%d", config.Get().App.Port)); err != nil && err != http.ErrServerClosed {
			echo.Logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := echo.Shutdown(ctx); err != nil {
		echo.Logger.Fatal(err)
	}
	shutdownFunc()

	log.Println("server stopped")
}
