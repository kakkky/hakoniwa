package usecase

type Usecases struct {
	RegisterResident        *RegisterResident
	RegisterBuildingManager *RegisterBuildingManager
	GetBuildingManager      *GetBuildingManager
	UpdateBuildingManager   *UpdateBuildingManager
	SendMessage             *SendMessageFromBuildingManagerToResident
}
