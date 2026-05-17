package domain

import (
	"context"
)

type AgentCommander interface {
	PublishCommand(ctx context.Context, cmd AgentCommand) error
}

type AgentCommand interface {
	_isAgentCommand()
}

type AgentCommandInbox chan AgentCommand

type PublishEventCommand struct {
	Event Event
}

type AddResidentAgentCommand struct {
	Resident Resident // snapshotを丸ごと
}

func (AddResidentAgentCommand) _isAgentCommand() {}
func (PublishEventCommand) _isAgentCommand()     {}
