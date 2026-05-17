package usecase

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type RegisterBuildingManager struct {
	repository domain.BuildingManagerRepository
}

func NewRegisterBuildingManager(
	repository domain.BuildingManagerRepository,
) *RegisterBuildingManager {
	return &RegisterBuildingManager{
		repository: repository,
	}
}

func (rbm *RegisterBuildingManager) Exec(ctx context.Context, buildingManager *domain.BuildingManager) error {
	if err := rbm.repository.Save(buildingManager); err != nil {
		return err
	}
	return nil
}
