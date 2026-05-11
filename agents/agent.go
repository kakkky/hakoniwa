package agents

import "github.com/kakkky/hakoniwa/domain"

type ID string

type base struct {
	ID          ID
	Inbox       eventInbox
	llmProvider domain.LLMProvider
}
