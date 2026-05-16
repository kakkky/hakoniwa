package agents

import (
	"github.com/kakkky/hakoniwa/domain"
)

type id string

type name string

type agentBase struct {
	id          id
	name        name
	inbox       agentEventInbox
	sendEvent   func(agentEvent)
	llmProvider domain.LLMProvider
	llmPrompt   *domain.LLMPrompts
}

func newAgentBase(
	sendEvent func(agentEvent),
	llmProvider domain.LLMProvider,
) *agentBase {
	inbox := make(chan agentEvent, 16)
	return &agentBase{
		inbox:       inbox,
		sendEvent:   sendEvent,
		llmProvider: llmProvider,
		llmPrompt:   &domain.LLMPrompts{},
	}
}
