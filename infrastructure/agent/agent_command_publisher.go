package agent

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type AgentCommandPublisher struct {
	agentCommandCh domain.AgentCommandCh
}

func NewAgentCommandPublisher(ch domain.AgentCommandCh) *AgentCommandPublisher {
	return &AgentCommandPublisher{
		agentCommandCh: ch,
	}
}

func (ac *AgentCommandPublisher) PublishCommand(ctx context.Context, cmd domain.AgentCommand) error {
	ac.agentCommandCh <- cmd
	return nil
}
