//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	agentpkg "github.com/kakkky/hakoniwa/agents"
	"github.com/kakkky/hakoniwa/config"
	"github.com/kakkky/hakoniwa/infrastructure"
)

func initializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		infrastructure.Set,
		agentpkg.Set,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
