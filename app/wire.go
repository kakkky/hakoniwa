//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	agentpkg "github.com/kakkky/hakoniwa/agents"
)

var Set = wire.NewSet(
	wire.Struct(new(App), "*"),
)
