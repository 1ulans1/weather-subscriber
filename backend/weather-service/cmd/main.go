package main

import (
	"go.uber.org/zap"

	"weather-service/config"
	"weather-service/internal/adapters/db"
	"weather-service/internal/adapters/external"
	"weather-service/internal/adapters/grpc"
	"weather-service/internal/core/services"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	gormDB := db.Gorm(logger)

	repo := db.NewWeatherRepo(gormDB)
	apiClient := external.NewWeatherAPIClient(config.Conf.Weather.Api.Key)
	weatherSvc := services.NewWeatherService(repo, apiClient)

	grpcPort := ":" + config.Conf.Port
	sugar.Infof("Starting WeatherService gRPC on %s", grpcPort)
	grpc.ServeGRPC(weatherSvc, grpcPort)
}
