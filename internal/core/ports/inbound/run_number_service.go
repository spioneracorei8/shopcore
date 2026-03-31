package inbound

import (
	"context"
	"shopcore/internal/core/domain"
)

type RunNumberUsecase interface {
	CreateRunNumber(ctx context.Context, rn *domain.RunNumber) error
	FetchRunNumber(ctx context.Context) (*domain.RunNumber, error)
	UpdateRunNumber(ctx context.Context) error
}
