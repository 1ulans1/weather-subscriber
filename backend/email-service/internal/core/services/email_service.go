package services

import (
	"email-service/internal/core/domain"
	"email-service/internal/core/ports"
	"fmt"
	"go.uber.org/zap"
)

type EmailService struct {
	sender  ports.EmailSender
	logger  *zap.Logger
	baseURL string
}

func NewEmailService(sender ports.EmailSender, logger *zap.Logger, baseURL string) *EmailService {
	return &EmailService{sender: sender, logger: logger, baseURL: baseURL}
}

type ProcessConfirmFunc func(msg *domain.ConfirmMessage) error

func (s *EmailService) ProcessConfirm(msg *domain.ConfirmMessage) error {
	s.logger.Info("Processing confirm email", zap.String("email", msg.Email))
	url := fmt.Sprintf("%s/api/confirm/%s", s.baseURL, msg.Token)
	if err := s.sender.SendConfirmation(msg.Email, url); err != nil {
		s.logger.Error("Failed to send confirmation email", zap.Error(err))
		return err
	}
	return nil
}

type ProcessNotifyFunc func(msg *domain.NotifyMessage) error

func (s *EmailService) ProcessNotify(msg *domain.NotifyMessage) error {
	s.logger.Info("Processing notification email", zap.String("email", msg.Email))
	unsubURL := fmt.Sprintf("%s/api/unsubscribe/%s", s.baseURL, msg.UnsubToken)
	if err := s.sender.SendNotification(msg.Email, msg.City, msg.Weather.Temperature, msg.Weather.Condition, unsubURL); err != nil {
		s.logger.Error("Failed to send notification email", zap.Error(err))
		return err
	}
	return nil
}
