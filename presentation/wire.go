//go:build wireinject
// +build wireinject

package presentation

import (
	"github.com/google/wire"
	"github.com/kakkky/hakoniwa/presentation/tui"
)

var Set = wire.NewSet(
	tui.NewTUI,
)
