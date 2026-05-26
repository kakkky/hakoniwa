package agents

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type commandSubscriber struct {
	agentCommandCh domain.AgentCommandCh
	routes         map[domain.AgentCommandKey]commandHandlerFunc
}

type commandHandlerFunc func(ctx context.Context, cmd domain.AgentCommand) error

func newCommandSubscriber(ch domain.AgentCommandCh) *commandSubscriber {
	return &commandSubscriber{
		agentCommandCh: ch,
		routes:         make(map[domain.AgentCommandKey]commandHandlerFunc),
	}
}

func registerHandler[T domain.AgentCommand](
	cs *commandSubscriber,
	cmd T,
	handler func(ctx context.Context, cmd T) error) {
	cs.routes[cmd.CommandKey()] = func(ctx context.Context, cmd domain.AgentCommand) error {
		return handler(ctx, cmd.(T))
	}
}

func (cs *commandSubscriber) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case cmd := <-cs.agentCommandCh:
			if handler, ok := cs.routes[cmd.CommandKey()]; ok {
				handler(ctx, cmd)
			}
		}
	}
}
