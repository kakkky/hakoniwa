package agents

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type commandSubscriber struct {
	cmdInbox domain.AgentCommandInbox
	routes   map[domain.AgentCommandKey]commandHandlerFunc
}

type commandHandlerFunc func(ctx context.Context, cmd domain.AgentCommand) error

func newCommandSubscriber() *commandSubscriber {
	return &commandSubscriber{
		cmdInbox: make(domain.AgentCommandInbox, 32),
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
		case cmd := <-cs.cmdInbox:
			if handler, ok := cs.routes[cmd.CommandKey()]; ok {
				handler(ctx, cmd)
			}
		}
	}
}
