package agents

import (
	"github.com/kakkky/hakoniwa/domain"
)

type agentBase struct {
	inbox       agentEventInbox
	sendEvent   func(domain.Event)
	llmProvider domain.LLMProvider
}

func newAgentBase(
	sendEvent func(domain.Event),
	llmProvider domain.LLMProvider,
) agentBase {
	inbox := make(chan domain.Event, 16)
	return agentBase{
		inbox:       inbox,
		sendEvent:   sendEvent,
		llmProvider: llmProvider,
	}
}
