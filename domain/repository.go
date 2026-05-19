package domain

//go:generate mockgen -source=repository.go -destination=../testing/mock/mock_repository.go -package=mock

// Reposutory Interfaces
type ResidentRepository interface {
	Save(resident *Resident) error
	SaveAll(residents []*Resident) error
	GetAll() ([]*Resident, error)
	// FindByID(id string) (*Resident, error)
	// DeleteByID(id string) error
}

type BuildingManagerRepository interface {
	Save(BuildingManager *BuildingManager) error
	Get() (BuildingManager, error)
}
