package server

import (
	"context"
	"fmt"
	healthcheckDeliveryGrpc "grpc-demo/internal/healthcheck/v2/delivery/grpc"
	healthcheck_usecase "grpc-demo/internal/healthcheck/v2/usecase"
	notificationDeliveryGrpc "grpc-demo/internal/notification/v2/delivery/grpc"
	notification_usecase "grpc-demo/internal/notification/v2/usecase"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RunGrpcServer ...
func RunGrpcServer(ctx context.Context, port string) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))

	if err != nil {
		fmt.Println(err)
		return err
	}

	options := []grpc.ServerOption{}
	server := grpc.NewServer(options...)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			fmt.Println("Shutting down GRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// register healthcheck usecase
	healthcheckUsecase := healthcheck_usecase.NewHealthCheckUsecase()
	healthcheckDeliveryGrpc.NewHealthCheckHandler(server, healthcheckUsecase)
	// register notification usecase
	notificationUsecase := notification_usecase.NewNotificationUsecase()
	nu := notificationDeliveryGrpc.NewNotificationUserHandler(server, notificationUsecase)

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		count := 0
		for range ticker.C {
			count++
			if count%2 == 0 {
				nu.SendUserMessage(2, fmt.Sprintf("[***SERVER***]Test - %d", time.Now().Unix()))
				continue
			}
			nu.SendUserMessage(1, fmt.Sprintf("[***SERVER***]Test - %d", time.Now().Unix()))
		}
	}()

	reflection.Register(server)
	fmt.Println("Starting GRPC server...")
	err = server.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return server.Serve(listen)
}
