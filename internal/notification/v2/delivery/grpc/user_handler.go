package grpc

import (
	"context"
	"fmt"
	"grpc-demo/internal/notification/v2/usecase"
	pb "grpc-demo/protobuf/notification/v2"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"google.golang.org/grpc"
)

type NotificationUserHandler struct {
	pb.UnimplementedNoticationServiceServer
	usecase       usecase.INotificationUseCase
	users         sync.Map
	usersLRUCache *lru.Cache
}

type userStream struct {
	stream pb.NoticationService_BidirectionStreamingServer
	close  chan<- bool
}

// NewNotificationUserHandler ...
func NewNotificationUserHandler(gserver *grpc.Server, uc usecase.INotificationUseCase) *NotificationUserHandler {
	newLRU, _ := lru.New(128)
	server := &NotificationUserHandler{usecase: uc, usersLRUCache: newLRU}
	pb.RegisterNoticationServiceServer(gserver, server)
	return server
}

func (n *NotificationUserHandler) SendMessage(ctx context.Context, req *pb.ClientRequest) (*pb.ServerResponse, error) {
	mess, err := n.usecase.SendMessage(ctx, req.Message)
	if err != nil {
		return nil, err
	}
	return &pb.ServerResponse{Status: 1, Message: mess}, nil
}

func (n *NotificationUserHandler) ClientStreaming(stream pb.NoticationService_ClientStreamingServer) error {
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

func (n *NotificationUserHandler) ServerStreaming(req *pb.ClientRequest, stream pb.NoticationService_ServerStreamingServer) error {
	mess, err := n.usecase.ServerStreaming(stream.Context(), req.Message)
	if err != nil {
		return err
	}

	stream.Send(&pb.ServerResponse{Status: 1, Message: mess})
	time.Sleep(500 * time.Microsecond)
	stream.Send(&pb.ServerResponse{Status: 1, Message: "Server streaming ended."})
	return nil
}

func (n *NotificationUserHandler) BidirectionStreaming(stream pb.NoticationService_BidirectionStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			continue
		}

		if err != nil {
			fmt.Println("BidirectionStreaming error: ", err.Error())
			return err
		}

		mess, err := n.usecase.BidirectionStreaming(stream.Context(), req.Message)
		if err != nil {
			fmt.Println("Usecase BidirectionStreaming error: ", err.Error())
			continue
		}

		close := make(chan bool)
		n.users.Store(req.Id, userStream{stream: stream, close: close})
		n.usersLRUCache.Add(req.Id, userStream{stream: stream, close: close})

		stream.Send(&pb.ServerResponse{Status: 1, Message: mess})
		stream.Send(&pb.ServerResponse{
			Status:  1,
			Message: fmt.Sprintf("Execute message \"%s\" with user id %d done.", req.Message, req.Id),
		})

		ctx := stream.Context()
		for {
			select {
			case <-close:
				log.Printf("Closing stream for client ID: %d", req.Id)
				return nil
			case <-ctx.Done():
				log.Printf("Client ID %d has disconnected", req.Id)
				n.users.Delete(req.Id)
				n.usersLRUCache.Remove(req.Id)
				return nil
			}
		}
	}
}

func (n *NotificationUserHandler) SendUserMessage(ID int64, mess string) error {
	fmt.Printf("ID: %d - Mess: %s\n", ID, mess)
	// map
	// s, ok := n.users.Load(ID)
	// lru
	s, ok := n.usersLRUCache.Get(ID)
	if !ok {
		fmt.Printf("Load user stream fail [%d]\n", ID)
		return nil
	}

	stream, ok := s.(userStream)
	if !ok {
		fmt.Println("User disconnect")
	}

	err := stream.stream.Send(&pb.ServerResponse{
		Status:  1,
		Message: fmt.Sprintf("Server send message %s to user id %d", mess, ID),
	})

	if err != nil {
		fmt.Println("Server send error: ", err.Error())
	}

	return nil
}
