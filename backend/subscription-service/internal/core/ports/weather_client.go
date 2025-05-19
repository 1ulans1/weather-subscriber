package ports

import "context"

type WeatherData struct {
	Temperature float64
	Condition   string
}

type WeatherClient interface {
	GetCurrentWeather(ctx context.Context, city string) (*WeatherData, error)
}
