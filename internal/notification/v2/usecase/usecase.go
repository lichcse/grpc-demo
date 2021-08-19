package usecase

import (
	"context"
)

type notificationUseCase struct{}

// NewNotificationUsecase ...
func NewNotificationUsecase() INotificationUseCase {
	return &notificationUseCase{}
}

func (s *notificationUseCase) SendMessage(ctx context.Context, message string) (string, error) {
	return "I received the message: " + message, nil
}

func (s *notificationUseCase) ClientStreaming(ctx context.Context, message string) (string, error) {
	return message, nil
}

func (s *notificationUseCase) ServerStreaming(ctx context.Context, message string) (string, error) {
	return message, nil
}

func (s *notificationUseCase) BidirectionStreaming(ctx context.Context, message string) (string, error) {
	return "I received the message: " + message, nil
}
