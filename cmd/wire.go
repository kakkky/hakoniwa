//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kakkky/hakoniwa/agents"
	"github.com/kakkky/hakoniwa/config"
	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/infrastructure"
	"github.com/kakkky/hakoniwa/infrastructure/agent"
	"github.com/kakkky/hakoniwa/presentation"
	"github.com/kakkky/hakoniwa/presentation/tui"
	"github.com/kakkky/hakoniwa/usecase"
)

type App struct {
	AgentRuntime *agents.Runtime
	UI           *tui.TUI
}

func initializeApp() (*App, error) {
	wire.Build(
		config.Set,
		infrastructure.Set,
		agents.Set,
		usecase.Set,
		presentation.Set,
		provideAgentCommander,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}

func provideAgentCommander(r *agents.Runtime) domain.AgentCommander {
	return agent.NewAgentCommander(r.CommandInbox())
}
