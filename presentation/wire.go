//go:build wireinject
// +build wireinject

package presentation

import (
	"github.com/google/wire"
	"github.com/kakkky/hakoniwa/presentation/ui"
)

var Set = wire.NewSet(
	ui.NewUI,
)
