package domain

//go:generate mockgen -source=llm_provider.go -destination=../testing/mock/mock_llm_provider.go -package=mock

import (
	"context"
	"encoding/json"
)

// LLM Provider Interface
type LLMPrompts struct {
	System string
	User   string
}

type LLMResponse string

type LLMProvider interface {
	Generate(ctx context.Context, prompts LLMPrompts, responseSchema json.RawMessage) (json.RawMessage, error)
}
