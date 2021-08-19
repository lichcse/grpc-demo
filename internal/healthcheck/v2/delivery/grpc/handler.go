package grpc

import (
	"context"
	"grpc-demo/internal/healthcheck/v2/usecase"
	pb "grpc-demo/protobuf/healthcheck/v2"

	"google.golang.org/grpc"
)

type healthCheckHandler struct {
	pb.UnimplementedHealthCheckServer
	usecase usecase.IHealthCheckUseCase
}

// NewHealthCheckHandler ...
func NewHealthCheckHandler(gserver *grpc.Server, uc usecase.IHealthCheckUseCase) {
	server := &healthCheckHandler{
		usecase: uc,
	}
	pb.RegisterHealthCheckServer(gserver, server)
}

func (h *healthCheckHandler) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	err := h.usecase.Check(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

func (h *healthCheckHandler) Watch(ctx *pb.HealthCheckRequest, req pb.HealthCheck_WatchServer) error {
	return nil
}
