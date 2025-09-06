package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	microservice "github.com/ikotiki/go-rest-api-service-subscriptions"
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/config"
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/server/http/handler"

	"net/http"

	middleware_logger "github.com/ikotiki/go-rest-api-service-subscriptions/internal/server/http/middleware/logger"
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/service"
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage/postgresql"
	"github.com/ikotiki/go-rest-api-service-subscriptions/logger"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"
)

// @title           Subscription API
// @version         1.0
// @description     Subscription API service.

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	cfg := config.MustLoad()

	log := logger.InitLogger(cfg.Env, true)
	log.Debug().Msg("logger initialized")

	log.Info().Interface("config", cfg).Str("version", "1.0").Msg("starting microservice")
	log.Debug().Msg("debug messages are enabled")
	log.Trace().Msg("trace messages are enabled")

	// Create DB instance
	pgdb, err := postgresql.NewSQLStorage(postgresql.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBname,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}
	log.Info().Msg("database connected")

	store := storage.Storage{
		Subscriptions: postgresql.NewSubscriptionsStore(pgdb),
	}
	srv := service.NewService(store)

	if !cfg.HTTPServer.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	middlewares := []gin.HandlerFunc{}

	if cfg.HTTPServer.Auth {
		log.Info().Msg("basic auth enabled")
		middlewares = append(middlewares, gin.BasicAuth(cfg.HTTPServer.GetUsers()))
	}

	group := router.Group("/api/v1", middlewares...)
	group.Use(middleware_logger.New(log))

	handlers := handler.New(srv)
	handlers.InitRoutes(group)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	server := microservice.NewServer(
		&http.Server{
			Addr:         cfg.HTTPServer.Addr,
			Handler:      router,
			ReadTimeout:  cfg.HTTPServer.Timeout,
			WriteTimeout: cfg.HTTPServer.Timeout,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		})

	log.Info().Msgf("starting server on host %s", server.Addr())

	go func() {
		if err := server.Run(); err != nil {
			log.Fatal().Err(err).Msg("error starting server")
		}
	}()

	<-quit

	log.Info().Msg("stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("can't gracefully shutdown server")
	}

	log.Info().Msg("server stopped")
}
