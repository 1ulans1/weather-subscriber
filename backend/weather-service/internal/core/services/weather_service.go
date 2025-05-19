package services

import (
	"context"
	"time"
	"weather-service/internal/core/domain"
	"weather-service/internal/core/ports"
)

type WeatherService struct {
	repo ports.WeatherRepo
	api  ports.WeatherAPIClient
	ttl  time.Duration
}

func NewWeatherService(repo ports.WeatherRepo, api ports.WeatherAPIClient) *WeatherService {
	return &WeatherService{repo: repo, api: api, ttl: 30 * time.Minute}
}

func (s *WeatherService) Get(ctx context.Context, location string) (*domain.Weather, error) {
	w, err := s.repo.Get(location)
	if err == nil {
		if time.Since(w.UpdatedAt) < s.ttl {
			return w, nil
		}
	}
	fresh, err := s.api.Fetch(ctx, location)
	if err != nil {
		return w, err
	}
	if err := s.repo.Save(fresh); err != nil {
		return nil, err
	}
	return fresh, nil
}
