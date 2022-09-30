package main

import (
	"context"
	"github.com/Feokrat/music-dating-app/notifications/internal/notifications"
	"github.com/Feokrat/music-dating-app/notifications/pkg/database"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Feokrat/music-dating-app/notifications/internal/config"
	"github.com/Feokrat/music-dating-app/notifications/pkg/HTTPserver"
	"github.com/gin-gonic/gin"
)

var configFile = "configs/config"

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)

	cfg, err := config.Init(configFile, logger)
	if err != nil {
		logger.Fatalf("failed to load application configuration: %s", err)
	}

	db, err := database.NewPostgresDB(cfg.Postgresql, logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer database.ClosePostgresDB(db)

	handlers := buildHandler(db, logger)
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

func buildHandler(db *sqlx.DB, logger *log.Logger) http.Handler {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	rg := router.Group("/api/v1")

	chatRepository := notifications.NewChatRepository(db, logger)
	messagesRepository := notifications.NewMessageRepository(db, logger)
	messagesStatusesRepository := notifications.NewMessageStatusesRepository(db, logger)
	service := notifications.NewChatService(logger, chatRepository, messagesRepository, messagesStatusesRepository)
	notifications.RegisterHandlers(rg, service, logger)

	return router
}
