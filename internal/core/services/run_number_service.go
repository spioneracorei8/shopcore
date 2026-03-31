package services

import (
	"context"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/inbound"
	"shopcore/internal/core/ports/outbound"
)

type runNumberUsecase struct {
	runNumberRepo outbound.RunNumberRepository
}

func NewRunNumberUsecaseImpl(runNumberRepo outbound.RunNumberRepository) inbound.RunNumberUsecase {
	return &runNumberUsecase{
		runNumberRepo: runNumberRepo,
	}
}

func (u *runNumberUsecase) CreateRunNumber(ctx context.Context, rn *domain.RunNumber) error {
	rn.GenObjectID()
	return u.runNumberRepo.CreateRunNumber(ctx, rn)
}

func (u *runNumberUsecase) FetchRunNumber(ctx context.Context) (*domain.RunNumber, error) {
	return u.runNumberRepo.FetchRunNumber(ctx)
}

func (u *runNumberUsecase) UpdateRunNumber(ctx context.Context) error {
	rn, err := u.FetchRunNumber(ctx)
	if err != nil {
		return err
	}
	rn.Running++
	return u.runNumberRepo.UpdateRunNumber(ctx, rn)
}
