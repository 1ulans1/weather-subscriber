package ports

import (
	"context"
	"subscription-service/internal/core/domain"
	"time"
)

type SubscriptionRepo interface {
	GetByEmail(ctx context.Context, email string) (*domain.Subscription, error)
	CreateSubscription(ctx context.Context, sub *domain.Subscription) error
	UpdateSubscription(ctx context.Context, sub *domain.Subscription) error
	GetSubscriptionByToken(ctx context.Context, token string) (*domain.Subscription, error)
	CreatePendingConfirmation(ctx context.Context, pc *domain.PendingConfirmation) error
	GetPendingConfirmationByToken(ctx context.Context, token string) (*domain.PendingConfirmation, error)
	DeletePendingConfirmation(ctx context.Context, id string) error
	FindDue(ctx context.Context, now time.Time) ([]*domain.Subscription, error)
}
