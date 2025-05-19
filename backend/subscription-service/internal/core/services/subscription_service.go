package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"subscription-service/internal/core/domain"
	"subscription-service/internal/core/ports"
)

type SubscriptionService struct {
	repo     ports.SubscriptionRepo
	emailPub ports.EmailPublisher
	weather  ports.WeatherClient
	logger   *zap.Logger
}

func NewSubscriptionService(r ports.SubscriptionRepo, e ports.EmailPublisher, w ports.WeatherClient, l *zap.Logger) *SubscriptionService {
	return &SubscriptionService{repo: r, emailPub: e, weather: w, logger: l}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, email, city, freq string) error {
	if email == "" || city == "" || (freq != "hourly" && freq != "daily") {
		return errors.New("invalid input")
	}
	s.logger.Info("Subscription request", zap.String("email", email), zap.String("city", city), zap.String("frequency", freq))
	token := uuid.NewString()
	pc := &domain.PendingConfirmation{
		ID:        uuid.NewString(),
		Email:     email,
		City:      city,
		Frequency: freq,
		Token:     token,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreatePendingConfirmation(ctx, pc); err != nil {
		s.logger.Error("Failed to create pending confirmation", zap.Error(err))
		return err
	}
	if err := s.emailPub.PublishConfirm(email, token); err != nil {
		s.logger.Error("Failed to publish confirmation email", zap.Error(err))
		return err
	}
	s.logger.Info("Confirmation email sent", zap.String("email", email))
	return nil
}

func (s *SubscriptionService) Confirm(ctx context.Context, token string) error {
	pc, err := s.repo.GetPendingConfirmationByToken(ctx, token)
	if err != nil {
		s.logger.Error("Invalid confirmation token", zap.Error(err))
		return errors.New("invalid or expired token")
	}
	s.logger.Info("Confirming subscription", zap.String("email", pc.Email))
	sub, err := s.repo.GetByEmail(ctx, pc.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to check existing subscription", zap.Error(err))
		return err
	}
	now := time.Now()
	if sub == nil {
		sub = &domain.Subscription{
			ID:          uuid.NewString(),
			Email:       pc.Email,
			City:        pc.City,
			Frequency:   pc.Frequency,
			Token:       uuid.NewString(),
			ConfirmedAt: &now,
			CreatedAt:   now,
		}
		if err := s.repo.CreateSubscription(ctx, sub); err != nil {
			s.logger.Error("Failed to create subscription", zap.Error(err))
			return err
		}
		s.logger.Info("New subscription created", zap.String("email", sub.Email))
	} else {
		if sub.UnsubscribedAt != nil {
			return errors.New("cannot update unsubscribed subscription")
		}
		sub.City = pc.City
		sub.Frequency = pc.Frequency
		sub.Token = uuid.NewString() // Regenerate token for security
		if err := s.repo.UpdateSubscription(ctx, sub); err != nil {
			s.logger.Error("Failed to update subscription", zap.Error(err))
			return err
		}
		s.logger.Info("Subscription updated", zap.String("email", sub.Email))
	}
	if err := s.repo.DeletePendingConfirmation(ctx, pc.ID); err != nil {
		s.logger.Warn("Failed to delete pending confirmation", zap.Error(err))
	}
	return nil
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, token string) error {
	sub, err := s.repo.GetSubscriptionByToken(ctx, token)
	if err != nil {
		s.logger.Error("Invalid unsubscribe token", zap.Error(err))
		return errors.New("invalid token")
	}
	if sub.UnsubscribedAt != nil {
		return errors.New("already unsubscribed")
	}
	s.logger.Info("Unsubscribing", zap.String("email", sub.Email))
	now := time.Now()
	sub.UnsubscribedAt = &now
	if err := s.repo.UpdateSubscription(ctx, sub); err != nil {
		s.logger.Error("Failed to unsubscribe", zap.Error(err))
		return err
	}
	s.logger.Info("Unsubscribed successfully", zap.String("email", sub.Email))
	return nil
}

func (s *SubscriptionService) GetWeather(ctx context.Context, city string) (*ports.WeatherData, error) {
	return s.weather.GetCurrentWeather(ctx, city)
}

func (s *SubscriptionService) NotifyDue(ctx context.Context) error {
	now := time.Now()
	subs, err := s.repo.FindDue(ctx, now)
	if err != nil {
		s.logger.Error("Failed to find due subscriptions", zap.Error(err))
		return err
	}
	for _, sub := range subs {
		wd, err := s.weather.GetCurrentWeather(ctx, sub.City)
		if err != nil {
			s.logger.Warn("Failed to get weather data", zap.String("city", sub.City), zap.Error(err))
			continue
		}
		temp := fmt.Sprintf("%.1f", wd.Temperature)
		if err := s.emailPub.PublishNotification(sub.Email, sub.City, temp, wd.Condition, sub.Token); err != nil {
			s.logger.Warn("Failed to send notification", zap.String("email", sub.Email), zap.Error(err))
			continue
		}
		sub.LastNotifiedAt = &now
		if err := s.repo.UpdateSubscription(ctx, sub); err != nil {
			s.logger.Warn("Failed to update last notified time", zap.String("email", sub.Email), zap.Error(err))
		}
		s.logger.Info("Notification sent", zap.String("email", sub.Email))
	}
	return nil
}
