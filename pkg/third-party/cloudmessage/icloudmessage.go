package cloudmessage

import "context"

// ICloudMessageService ...
type ICloudMessageService interface {
	SendMessage(ctx context.Context, message string) (string, error)
}
