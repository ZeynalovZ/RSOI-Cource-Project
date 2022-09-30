package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Feokrat/music-dating-app/users/internal/config"
	"github.com/Feokrat/music-dating-app/users/internal/music"
	"github.com/Feokrat/music-dating-app/users/internal/user"
	"github.com/Feokrat/music-dating-app/users/pkg/HTTPserver"
	"github.com/Feokrat/music-dating-app/users/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

	userRepository := user.NewRepository(db, logger)
	userService := user.NewService(userRepository, logger)
	user.RegisterHandlers(rg.Group("/users"), userService, logger)

	musicRepository := music.NewRepository(db, logger)
	musicService := music.NewService(musicRepository, logger)
	music.RegisterHandlers(rg.Group("/musics"), musicService, logger)

	return router
}
