package usecase

type Usecases struct {
	RegisterResident        *RegisterResident
	RegisterBuildingManager *RegisterBuildingManager
	GetBuildingManager      *GetBuildingManager
	SendMessage             *SendMessageFromBuildingManagerToResident
}
