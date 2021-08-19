package main

import (
	"context"
	pb "grpc-demo/protobuf/notification/v2"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Println("Please enter one of the param: SendMessage, ClientStreaming, ServerStreaming, BidirectionStreaming")
	}

	clientConn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection error %v", err)
	}
	defer clientConn.Close()

	client := pb.NewNoticationServiceClient(clientConn)
	for _, m := range args {
		switch m {
		case "SendMessage":
			SendMessage(client)
		case "ClientStreaming":
			ClientStreaming(client)
		case "ServerStreaming":
			ServerStreaming(client)
		case "BidirectionStreaming":
			BidirectionStreaming(client)
		default:
			log.Println("Args invalid (Accept: SendMessage, ClientStreaming, ServerStreaming, BidirectionStreaming)")
		}
	}
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
	stream, err := client.BidirectionStreaming(context.Background())
	if err != nil {
		log.Fatalf("Bidirection streaming error %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Client streaming send message
	go func(stream pb.NoticationService_BidirectionStreamingClient, wg *sync.WaitGroup) {
		defer wg.Done()
		mess := []string{"Have", "a", "good", "one!"}
		for _, m := range mess {
			err := stream.Send(&pb.ClientRequest{Message: m})
			if err != nil {
				log.Fatalf("Bidirection streaming send message error %v", err)
			}

			log.Println("[BidirectionStreaming] client message: ", m)
			time.Sleep(1000 * time.Microsecond)
		}
		stream.CloseSend()
	}(stream, &wg)

	// Server streaming response message
	go func(stream pb.NoticationService_BidirectionStreamingClient, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("[BidirectionStreaming] done!!!")
				break
			}

			if err != nil {
				log.Fatalf("Server streaming error %v", err)
			}

			log.Println("[BidirectionStreaming] Server response: ", res.Message)
		}
	}(stream, &wg)

	wg.Wait()
}
