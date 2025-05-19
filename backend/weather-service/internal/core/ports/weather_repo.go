package ports

import "weather-service/internal/core/domain"

type WeatherRepo interface {
	Save(weather *domain.Weather) error
	Get(location string) (*domain.Weather, error)
}
