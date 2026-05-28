package usecase

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type GetBuildingManager struct {
	repository domain.BuildingManagerRepository
}

func NewGetBuildingManager(
	repository domain.BuildingManagerRepository,
) *GetBuildingManager {
	return &GetBuildingManager{
		repository: repository,
	}
}

func (gbm *GetBuildingManager) Exec(ctx context.Context) (*domain.BuildingManager, error) {
	bm, err := gbm.repository.Get()
	if err != nil {
		return nil, err
	}
	return bm, nil
}
