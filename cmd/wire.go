//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kakkky/hakoniwa/app"
	agentpkg "github.com/kakkky/hakoniwa/agents"
	"github.com/kakkky/hakoniwa/config"
	"github.com/kakkky/hakoniwa/infrastructure"
)

func InitializeApp(cfg *config.Config) (*app.App, error) {
	wire.Build(
		infrastructure.Set,
		agentpkg.Set,
		app.Set,
	)
	return nil, nil
}
