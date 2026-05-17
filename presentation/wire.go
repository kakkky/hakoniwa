//go:build wireinject
// +build wireinject

package presentation

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewTextUserInterface,
)
