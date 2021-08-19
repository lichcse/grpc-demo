package grpc

import (
	"context"
	"fmt"
	"grpc-demo/internal/notification/v2/usecase"
	pb "grpc-demo/protobuf/notification/v2"
	"io"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type notificationHandler struct {
	pb.UnimplementedNoticationServiceServer
	usecase usecase.INotificationUseCase
}

// NewNotificationHandler ...
func NewNotificationHandler(gserver *grpc.Server, uc usecase.INotificationUseCase) {
	server := &notificationHandler{
		usecase: uc,
	}
	pb.RegisterNoticationServiceServer(gserver, server)
}

func (n *notificationHandler) SendMessage(ctx context.Context, req *pb.ClientRequest) (*pb.ServerResponse, error) {
	mess, err := n.usecase.SendMessage(ctx, req.Message)
	if err != nil {
		return nil, err
	}
	return &pb.ServerResponse{Status: 1, Message: mess}, nil
}

func (n *notificationHandler) ClientStreaming(stream pb.NoticationService_ClientStreamingServer) error {
	var mess []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.ServerResponse{
				Status:  1,
				Message: strings.Join(mess, " "),
			})
		}

		if err != nil {
			return err
		}

		res, err := n.usecase.ClientStreaming(stream.Context(), req.Message)
		if err != nil {
			return err
		}
		mess = append(mess, res)
	}
}

func (n *notificationHandler) ServerStreaming(req *pb.ClientRequest, stream pb.NoticationService_ServerStreamingServer) error {
	mess, err := n.usecase.ServerStreaming(stream.Context(), req.Message)
	if err != nil {
		return err
	}

	stream.Send(&pb.ServerResponse{Status: 1, Message: mess})
	time.Sleep(500 * time.Microsecond)
	stream.Send(&pb.ServerResponse{Status: 1, Message: "Server streaming ended."})
	return nil
}

func (n *notificationHandler) BidirectionStreaming(stream pb.NoticationService_BidirectionStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		mess, err := n.usecase.BidirectionStreaming(stream.Context(), req.Message)
		if err != nil {
			return err
		}

		stream.Send(&pb.ServerResponse{Status: 1, Message: mess})
		stream.Send(&pb.ServerResponse{
			Status:  1,
			Message: fmt.Sprintf("Execute message \"%s\" done.", req.Message),
		})
	}
}
