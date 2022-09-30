package main

import (
	"context"
	"github.com/Feokrat/music-dating-app/gateway/internal/TokenValidator"
	"github.com/Feokrat/music-dating-app/gateway/internal/gateway"
	"github.com/Feokrat/music-dating-app/gateway/internal/notifications"
	"github.com/Feokrat/music-dating-app/gateway/internal/session"
	"github.com/gin-contrib/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Feokrat/music-dating-app/gateway/internal/config"
	"github.com/Feokrat/music-dating-app/gateway/pkg/HTTPserver"
	"github.com/gin-gonic/gin"
)

var configFile = "configs/config"

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)

	cfg, err := config.Init(configFile, logger)
	if err != nil {
		logger.Fatalf("failed to load application configuration: %s", err)
	}

	handlers := buildHandler(cfg, logger)
	server := HTTPserver.NewHTTPserver(cfg, handlers)

	go func() {
		if err := server.Run(); err != nil {
			logger.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	sig := <-c
	logger.Println("Got signal:", sig)
	logger.Print("Shutting down server")

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	server.Stop(ctx)
}

func buildHandler(cfg *config.Config, logger *log.Logger) http.Handler {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.Default(),
		CORSMiddleware(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	rg := router.Group("/api/v1")
	gateway.RegisterUsersHandlers(rg.Group(""), gateway.NewUsersService(cfg.Services, logger),
		TokenValidator.NewValidationService(logger, cfg.Services), logger)

	session.RegisterAuthHandlers(rg.Group("/sessions"), session.NewSessionService(logger, cfg.Services),
		logger, gateway.NewUsersService(cfg.Services, logger))

	notifications.RegisterChatHandlers(rg.Group("/chats"), notifications.NewNotificationService(cfg.Services, logger),
		TokenValidator.NewValidationService(logger, cfg.Services), gateway.NewUsersService(cfg.Services, logger), logger)

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
