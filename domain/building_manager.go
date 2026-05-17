package domain

import "time"

type BuildManagerName string

type BuildingManager struct {
	Name        BuildManagerName
	Age         int
	AppointedAt time.Time
}

func NewBuildingManager(name string, age int, now time.Time) *BuildingManager {
	return &BuildingManager{
		Name:        BuildManagerName(name),
		Age:         age,
		AppointedAt: now,
	}
}
