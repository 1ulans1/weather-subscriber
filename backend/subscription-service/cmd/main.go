package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"subscription-service/config"
	"subscription-service/internal/adapters/db"
	grpcAdapter "subscription-service/internal/adapters/grpc"
	httpAdapter "subscription-service/internal/adapters/http"
	"subscription-service/internal/adapters/messaging"
	"subscription-service/internal/core/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Couldn't set up the logger: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	gormDB := db.Gorm(logger)
	subRepo := db.NewSubscriptionRepo(gormDB)

	rabbitmqConnStr := fmt.Sprintf("amqp://%s:%s@%s/", config.Conf.Rabbit.User, config.Conf.Rabbit.Password, config.Conf.Rabbit.Host)
	emailPub, err := messaging.NewRabbitMQPublisher(rabbitmqConnStr)
	if err != nil {
		sugar.Fatalf("Can't connect to RabbitMQ: %v", err)
	}

	weatherClient, err := grpcAdapter.NewWeatherClient(config.Conf.Weather.GRPC.Host + ":" + config.Conf.Weather.GRPC.Port)
	if err != nil {
		sugar.Fatalf("Can't connect to weather service: %v", err)
	}

	subSvc := services.NewSubscriptionService(subRepo, emailPub, weatherClient, logger)

	app := fiber.New()
	handler := httpAdapter.NewHandler(subSvc)
	handler.RegisterRoutes(app)

	go func() {
		interval, _ := time.ParseDuration(config.Conf.Notification.Interval)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := subSvc.NotifyDue(context.Background()); err != nil {
					sugar.Errorf("Problem sending notifications: %v", err)
				}
			}
		}
	}()

	port := config.Conf.HTTP.Port
	sugar.Infof("Subscription Service is running on port %s", port)
	go func() {
		if err := app.Listen(":" + port); err != nil {
			sugar.Fatalf("Can't start the server: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	sugar.Info("Turning off the subscription service...")
	app.Shutdown()
}
