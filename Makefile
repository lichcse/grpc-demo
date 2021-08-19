server:
	go run cmd/api/main.go

client:
	go run cmd/client/notification/main.go SendMessage ClientStreaming ServerStreaming BidirectionStreaming