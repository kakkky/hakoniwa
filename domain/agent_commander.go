package domain

//go:generate mockgen -source=agent_commander.go -destination=mock/mock_agent_commander.go -package=mock

import (
	"context"
)

type AgentCommander interface {
	PublishCommand(ctx context.Context, cmd AgentCommand) error
}

type AgentCommand interface {
	_isAgentCommand()
	CommandKey() AgentCommandKey
}

type AgentCommandKey string

const (
	PublishEvent     AgentCommandKey = "publish_event"
	AddResidentAgent AgentCommandKey = "add_resident_agent"
)

type AgentCommandInbox chan AgentCommand

type PublishEventCommand struct {
	Event Event
}

func (pe PublishEventCommand) CommandKey() AgentCommandKey {
	return PublishEvent
}

type AddResidentAgentCommand struct {
	Resident Resident // snapshotを丸ごと
}

func (ar AddResidentAgentCommand) CommandKey() AgentCommandKey {
	return AddResidentAgent
}

func (AddResidentAgentCommand) _isAgentCommand() {}
func (PublishEventCommand) _isAgentCommand()     {}
