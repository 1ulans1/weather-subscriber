package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"

	"email-service/config"
	"email-service/internal/adapters/email"
	"email-service/internal/adapters/messaging"
	"email-service/internal/core/services"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	rab := config.Conf.Rabbit
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s/", rab.User, rab.Password, rab.Host)

	smtpCfg := config.Conf.Email
	sender := email.NewSMTPSender(
		smtpCfg.SMTPHost,
		smtpCfg.SMTPUser,
		smtpCfg.SMTPPassword,
		smtpCfg.FromAddress,
	)

	emailSvc := services.NewEmailService(sender, logger, config.Conf.HTTP.BaseURL)

	consumer, err := messaging.NewRabbitMQConsumer(amqpURL)
	if err != nil {
		logger.Fatal("RabbitMQ connection failed", zap.Error(err))
	}

	go consumer.ConsumeConfirm(context.Background(), emailSvc.ProcessConfirm)
	go consumer.ConsumeNotification(context.Background(), emailSvc.ProcessNotify)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down Email Service...")
}
