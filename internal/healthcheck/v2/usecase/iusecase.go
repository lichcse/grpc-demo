package usecase

import "context"

// IHealthCheckUseCase ...
type IHealthCheckUseCase interface {
	Check(ctx context.Context) error
}
