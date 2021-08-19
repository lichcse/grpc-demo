package usecase

import "context"

// INotificationUseCase ...
type INotificationUseCase interface {
	SendMessage(ctx context.Context, message string) (string, error)
	ClientStreaming(ctx context.Context, message string) (string, error)
	ServerStreaming(ctx context.Context, message string) (string, error)
	BidirectionStreaming(ctx context.Context, message string) (string, error)
}
