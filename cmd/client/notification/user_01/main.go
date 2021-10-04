package main

import (
	"context"
	"fmt"
	pb "grpc-demo/protobuf/notification/v2"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	clientConn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection error %v", err)
	}
	defer clientConn.Close()

	client := pb.NewNoticationServiceClient(clientConn)
	BidirectionStreaming(client)
}

func SendMessage(client pb.NoticationServiceClient) {
	res, err := client.SendMessage(context.Background(), &pb.ClientRequest{Message: "Hello"})
	if err != nil {
		log.Fatalf("Send message error %v", err)
	}

	log.Println("[SendMessage] Server response: ", res.Message)
	log.Println("[SendMessage] done!!!")
}

func ClientStreaming(client pb.NoticationServiceClient) {
	stream, err := client.ClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("[Start] Client streaming message error %v", err)
	}

	mess := []string{"Have", "a", "good", "one!"}
	for _, m := range mess {
		err := stream.Send(&pb.ClientRequest{Message: m})
		if err != nil {
			log.Fatalf("Send message to server error %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("[Close] Client streaming message error %v", err)
	}

	log.Println("[ClientStreaming] Server response: ", res.Message)
	log.Println("[ClientStreaming] done!!!")
}

func ServerStreaming(client pb.NoticationServiceClient) {
	stream, err := client.ServerStreaming(context.Background(), &pb.ClientRequest{Message: "Have a good one!"})
	if err != nil {
		log.Fatalf("Server streaming error %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("[ServerStreaming] done!!!")
			return
		}

		log.Println("[ServerStreaming] Server response: ", res.Message)
	}
}

func BidirectionStreaming(client pb.NoticationServiceClient) {
	keepConnect := make(chan bool)
	stream, err := client.BidirectionStreaming(context.Background())
	if err != nil {
		log.Fatalf("Bidirection streaming error %v", err)
	}

	// Client streaming send message
	go func(stream pb.NoticationService_BidirectionStreamingClient) {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			mess := fmt.Sprintf("Schedule send message to user 1 - Time: %d", time.Now().Unix())
			err := stream.Send(&pb.ClientRequest{Id: 1, Message: mess})
			if err != nil {
				log.Fatalf("Bidirection streaming send message error %v", err)
			}

			log.Println("[BidirectionStreaming] client message: ", mess)
		}
	}(stream)

	// Server streaming response message
	go func(stream pb.NoticationService_BidirectionStreamingClient) {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				// log.Println("[BidirectionStreaming] done!!!")
				// break
				continue
			}

			if err != nil {
				log.Fatalf("Server streaming error %v", err)
			}

			log.Println("[BidirectionStreaming] Server response: ", res.Message)
		}
	}(stream)

	<-keepConnect
}
