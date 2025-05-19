package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
	"weather-service/internal/core/services"
	"weather-service/pb"
)

type server struct {
	svc *services.WeatherService
	pb.UnimplementedWeatherServiceServer
}

func NewGRPCServer(svc *services.WeatherService) *grpc.Server {
	grpcSrv := grpc.NewServer()
	pb.RegisterWeatherServiceServer(grpcSrv, &server{svc: svc})
	return grpcSrv
}

func (s *server) GetCurrentWeather(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	w, err := s.svc.Get(ctx, req.Location)
	if err != nil {
		return nil, err
	}
	return &pb.WeatherResponse{
		Location:    w.Location,
		Temperature: w.Temperature,
		Condition:   w.Condition,
		UpdatedAt:   w.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func ServeGRPC(svc *services.WeatherService, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcSrv := NewGRPCServer(svc)
	log.Printf("gRPC server listening on %s", port)
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
