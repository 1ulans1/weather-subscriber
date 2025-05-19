package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"subscription-service/internal/core/ports"
	"subscription-service/pb"
)

type WeatherClient struct {
	client pb.WeatherServiceClient
}

func NewWeatherClient(addr string) (ports.WeatherClient, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(addr, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &WeatherClient{client: pb.NewWeatherServiceClient(conn)}, nil
}

func (w *WeatherClient) GetCurrentWeather(ctx context.Context, city string) (*ports.WeatherData, error) {
	resp, err := w.client.GetCurrentWeather(ctx, &pb.WeatherRequest{Location: city})
	if err != nil {
		return nil, err
	}
	return &ports.WeatherData{
		Temperature: resp.Temperature,
		Condition:   resp.Condition,
	}, nil
}
