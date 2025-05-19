package domain

import "time"

type Subscription struct {
	ID             string
	Email          string
	City           string
	Frequency      string
	Token          string
	ConfirmedAt    *time.Time
	UnsubscribedAt *time.Time
	LastNotifiedAt *time.Time
	CreatedAt      time.Time
}
