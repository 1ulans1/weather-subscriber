package domain

import "time"

type PendingConfirmation struct {
	ID        string
	Email     string
	City      string
	Frequency string
	Token     string
	CreatedAt time.Time
}
