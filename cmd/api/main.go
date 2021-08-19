package main

import (
	"context"
	"fmt"
	"grpc-demo/cmd/api/server"
	"os"
)

// Config is configuration for Server
type Config struct {
	GRPCPort string
	HTTPPort string
}

// Server ...
type Server struct {
	config Config
}

func (s *Server) RunServer() error {
	ctx := context.Background()

	s.config.GRPCPort = "5001"
	s.config.HTTPPort = "9001"

	go func() {
		_ = server.RunHttpServer(ctx, s.config.GRPCPort, s.config.HTTPPort)
	}()

	return server.RunGrpcServer(ctx, s.config.GRPCPort)
}

func main() {
	srv := &Server{}
	if err := srv.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
