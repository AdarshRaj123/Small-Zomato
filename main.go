package main

import (
	"SmallZomato/database"
	_ "SmallZomato/docs"
	"SmallZomato/server"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutDownTimeOut = 10 * time.Second

// @title Small Zomato
// @version 1.0
// @description This is a sample server of Small Zomato .
// @schemes http
// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key
func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// create server instance
	srv := server.SetupRoutes()
	if err := database.ConnectAndMigrate(
		"localhost",
		"5432",
		"small_zomato",
		"local",
		"local",
		database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Print("migration successful!!")
	srv.Router.Get("/swagger/*", httpSwagger.Handler())

	go func() {
		if err := srv.Run(":3000"); err != nil && err != http.ErrServerClosed {
			logrus.Panicf("Failed to run server with error: %+v", err)
		}
	}()
	logrus.Print("Server started at :8080")

	<-done

	logrus.Info("shutting down server")
	if err := database.ShutdownDatabase(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}
	if err := srv.Shutdown(shutDownTimeOut); err != nil {
		logrus.WithError(err).Panic("failed to gracefully shutdown server")
	}
}
