package server

import (
	"context"
	"fmt"
	healthcheck_pb "grpc-demo/protobuf/healthcheck/v2"
	notification_pb "grpc-demo/protobuf/notification/v2"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// RunRestServer runs HTTP/REST gateway
func RunHttpServer(ctx context.Context, grpcPort string, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	fmt.Println("GRPC Port: ", grpcPort)
	port := fmt.Sprintf(":%s", grpcPort)

	err := healthcheck_pb.RegisterHealthCheckHandlerFromEndpoint(ctx, mux, port, opts)
	if err != nil {
		panic(fmt.Sprintf("Failed to start HTTP/REST gateway: %s", err.Error()))
	}

	err = notification_pb.RegisterNoticationServiceHandlerFromEndpoint(ctx, mux, port, opts)
	if err != nil {
		panic(fmt.Sprintf("Failed to start HTTP/REST gateway: %s", err.Error()))
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			fmt.Println("Shutting down HTTP server...")
		}
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	fmt.Println("Starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}
