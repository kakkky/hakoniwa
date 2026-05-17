package agent

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type AgentCommander struct {
	inbox domain.AgentCommandInbox
}

func NewAgentCommander(inbox domain.AgentCommandInbox) *AgentCommander {
	return &AgentCommander{
		inbox: inbox,
	}
}

func (ac *AgentCommander) PublishCommand(ctx context.Context, cmd domain.AgentCommand) error {
	ac.inbox <- cmd
	return nil
}
