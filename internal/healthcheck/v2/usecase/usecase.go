package usecase

import (
	"context"
)

type healthCheckUseCase struct{}

// NewHealthCheckUsecase ...
func NewHealthCheckUsecase() IHealthCheckUseCase {
	return &healthCheckUseCase{}
}

// Check ...
func (s *healthCheckUseCase) Check(ctx context.Context) error {
	return nil
}
