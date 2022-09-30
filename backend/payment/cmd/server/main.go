package main

import (
	"context"
	"github.com/Feokrat/music-dating-app/payment/internal/payments"
	"github.com/Feokrat/music-dating-app/payment/internal/payments/repositories"
	"github.com/Feokrat/music-dating-app/payment/pkg/database"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Feokrat/music-dating-app/payment/internal/config"
	"github.com/Feokrat/music-dating-app/payment/pkg/HTTPserver"
	"github.com/gin-gonic/gin"
)

var configFile = "configs/config"

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)

	cfg, err := config.Init(configFile, logger)
	if err != nil {
		logger.Fatalf("failed to load application configuration: %s", err)
	}

	db, err := database.NewPostgresDB(cfg.PostgreSQL, logger)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	handlers := buildHandler(logger, db, cfg)
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

func buildHandler(logger *log.Logger, db *sqlx.DB, cfg *config.Config) http.Handler {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	rg := router.Group("/payments")
	paymentsRepository := repositories.NewPaymentsRepository(db, logger)
	paymentsService := payments.NewPaymentService(logger, paymentsRepository)

	payments.RegisterHandlers(rg, paymentsService, logger)

	return router
}
