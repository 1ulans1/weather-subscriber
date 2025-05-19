package domain

import "time"

type Weather struct {
	Location    string
	Temperature float64
	Humidity    int
	Condition   string
	UpdatedAt   time.Time
}
