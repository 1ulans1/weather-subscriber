package db

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"subscription-service/internal/core/domain"
	"subscription-service/internal/core/ports"
)

type GormSubscription struct {
	ID             string `gorm:"primaryKey"`
	Email          string `gorm:"uniqueIndex;not null"`
	City           string `gorm:"not null"`
	Frequency      string `gorm:"not null"`
	Token          string `gorm:"uniqueIndex;not null"`
	ConfirmedAt    *time.Time
	UnsubscribedAt *time.Time
	LastNotifiedAt *time.Time
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

type GormPendingConfirmation struct {
	ID        string    `gorm:"primaryKey"`
	Email     string    `gorm:"not null"`
	City      string    `gorm:"not null"`
	Frequency string    `gorm:"not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) ports.SubscriptionRepo {
	repo := &subscriptionRepo{db: db}
	repo.init()
	return repo
}

func (r *subscriptionRepo) init() {
	if err := r.db.AutoMigrate(&GormSubscription{}, &GormPendingConfirmation{}); err != nil {
		panic(err)
	}
}

func (r *subscriptionRepo) GetByEmail(ctx context.Context, email string) (*domain.Subscription, error) {
	var rec GormSubscription
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainSubscription(&rec), nil
}

func (r *subscriptionRepo) CreateSubscription(ctx context.Context, sub *domain.Subscription) error {
	rec := GormSubscription{
		ID:          sub.ID,
		Email:       sub.Email,
		City:        sub.City,
		Frequency:   sub.Frequency,
		Token:       sub.Token,
		ConfirmedAt: sub.ConfirmedAt,
		CreatedAt:   sub.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&rec).Error
}

func (r *subscriptionRepo) UpdateSubscription(ctx context.Context, sub *domain.Subscription) error {
	updates := map[string]interface{}{
		"city":             sub.City,
		"frequency":        sub.Frequency,
		"token":            sub.Token,
		"unsubscribed_at":  sub.UnsubscribedAt,
		"last_notified_at": sub.LastNotifiedAt,
	}
	return r.db.WithContext(ctx).Model(&GormSubscription{}).Where("id = ?", sub.ID).Updates(updates).Error
}

func (r *subscriptionRepo) GetSubscriptionByToken(ctx context.Context, token string) (*domain.Subscription, error) {
	var rec GormSubscription
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&rec).Error; err != nil {
		return nil, err
	}
	return r.toDomainSubscription(&rec), nil
}

func (r *subscriptionRepo) CreatePendingConfirmation(ctx context.Context, pc *domain.PendingConfirmation) error {
	rec := GormPendingConfirmation{
		ID:        pc.ID,
		Email:     pc.Email,
		City:      pc.City,
		Frequency: pc.Frequency,
		Token:     pc.Token,
		CreatedAt: pc.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&rec).Error
}

func (r *subscriptionRepo) GetPendingConfirmationByToken(ctx context.Context, token string) (*domain.PendingConfirmation, error) {
	var rec GormPendingConfirmation
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&rec).Error; err != nil {
		return nil, err
	}
	return &domain.PendingConfirmation{
		ID:        rec.ID,
		Email:     rec.Email,
		City:      rec.City,
		Frequency: rec.Frequency,
		Token:     rec.Token,
		CreatedAt: rec.CreatedAt,
	}, nil
}

func (r *subscriptionRepo) DeletePendingConfirmation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&GormPendingConfirmation{}, "id = ?", id).Error
}

func (r *subscriptionRepo) FindDue(ctx context.Context, now time.Time) ([]*domain.Subscription, error) {
	var recs []GormSubscription
	hourlyCutoff := now.Add(-1 * time.Hour)
	dailyCutoff := now.Add(-24 * time.Hour)
	err := r.db.WithContext(ctx).
		Where("confirmed_at IS NOT NULL").
		Where("unsubscribed_at IS NULL").
		Where("(frequency = ? AND (last_notified_at IS NULL OR last_notified_at <= ?)) OR (frequency = ? AND (last_notified_at IS NULL OR last_notified_at <= ?))",
			"hourly", hourlyCutoff, "daily", dailyCutoff).
		Find(&recs).Error
	if err != nil {
		return nil, err
	}
	var subs []*domain.Subscription
	for _, rec := range recs {
		subs = append(subs, r.toDomainSubscription(&rec))
	}
	return subs, nil
}

func (r *subscriptionRepo) toDomainSubscription(rec *GormSubscription) *domain.Subscription {
	return &domain.Subscription{
		ID:             rec.ID,
		Email:          rec.Email,
		City:           rec.City,
		Frequency:      rec.Frequency,
		Token:          rec.Token,
		ConfirmedAt:    rec.ConfirmedAt,
		UnsubscribedAt: rec.UnsubscribedAt,
		LastNotifiedAt: rec.LastNotifiedAt,
		CreatedAt:      rec.CreatedAt,
	}
}
