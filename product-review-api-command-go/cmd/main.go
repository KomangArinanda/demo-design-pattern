package main

import (
	"context"
	"errors"
	httpRouter "example/product-review-api-command-go/api/http-router"
	"example/product-review-api-command-go/config"
	"example/product-review-api-command-go/internal/client"
	"example/product-review-api-command-go/internal/handler"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	product_review "example/product-review-api-command-go/internal/usecase/product_review"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	repository := repo.NewProductReviewRepo(cfg.DatabaseSimulatedLatency)
	repository.Seed()

	orderClient := client.NewOrderClient()
	productClient := client.NewProductClient()
	sharedService := service.NewProductReviewService()
	usecases := product_review.NewProductReviewUsecases(repository, orderClient, productClient, sharedService)
	httpHandler := handler.NewHandler(usecases)
	router := httpRouter.InitializeRouter(httpHandler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		log.Printf("product-review-api-command-go listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}
