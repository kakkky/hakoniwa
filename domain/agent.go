package domain

import (
	"context"
	"time"
)

type AgentCommand interface {
	_isAgentCommand()
}

type PublishTickCommand struct {
	Now   time.Time // worldTime 入れるなら別途
	Event WorldEvent
}

type AddResidentAgentCommand struct {
	Resident Resident // snapshotを丸ごと
}

type AgentCommandInbox chan AgentCommand

type AgentCommander interface {
	PublishTickEvent(ctx context.Context, cmd PublishTickCommand) error
	AddResidentAgent(ctx context.Context, cmd AddResidentAgentCommand) error
}

func (AddResidentAgentCommand) _isAgentCommand() {}
func (PublishTickCommand) _isAgentCommand()      {}

type AgentEventSubscriber interface {
	Subscribe()
}
