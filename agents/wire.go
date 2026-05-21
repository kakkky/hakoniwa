//go:build wireinject
// +build wireinject

package agents

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewRuntime,
	 wire.Struct(new(ToolKit), "*")
)
