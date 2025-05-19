package ports

import (
	"context"
	"weather-service/internal/core/domain"
)

type WeatherAPIClient interface {
	Fetch(ctx context.Context, location string) (*domain.Weather, error)
}
