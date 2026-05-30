package usecase

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type UpdateBuildingManager struct {
	repository domain.BuildingManagerRepository
}

func NewUpdateBuildingManager(
	repository domain.BuildingManagerRepository,
) *UpdateBuildingManager {
	return &UpdateBuildingManager{
		repository: repository,
	}
}

func (ubm *UpdateBuildingManager) Exec(ctx context.Context, bdm *domain.BuildingManager) error {
	originBdm, err := ubm.repository.Get()
	if err != nil {
		return err
	}
	bdm.AppointedAt = originBdm.AppointedAt

	if err := ubm.repository.Save(bdm); err != nil {
		return err
	}
	return nil
}
