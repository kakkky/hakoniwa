//go:build wireinject
// +build wireinject

package usecase

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewRegisterResident,
	NewRegisterBuildingManager,
	NewGetBuildingManager,
	NewSendMessageFromBuildingManagerToResident,
	wire.Struct(new(Usecases), "*"),
)
