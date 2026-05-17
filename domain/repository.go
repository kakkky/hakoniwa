package domain

// Reposutory Interfaces
type ResidentRepository interface {
	Save(resident *Resident) error
	SaveAll(residents []*Resident) error
	// FindByID(id string) (*Resident, error)
	// DeleteByID(id string) error
}

type BuildingManagerRepository interface {
	Save(BuildingManager *BuildingManager) error
	Get() (BuildingManager, error)
}
