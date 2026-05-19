package agent

import "github.com/kakkky/hakoniwa/domain"

func NewAgentCommandCh() domain.AgentCommandCh {
	return make(chan domain.AgentCommand, 32)
}
