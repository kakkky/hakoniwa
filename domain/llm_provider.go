package domain

//go:generate mockgen -source=llm_provider.go -destination=../testing/mock/mock_llm_provider.go -package=mock

import (
	"context"
	"encoding/json"
	"strings"
)

// LLM Provider Interface
type LLMPrompts struct {
	System strings.Builder
	User   strings.Builder
}

func (l *LLMPrompts) AddSystemPrompt(new string) {
	l.System.WriteString(new + "\n")
}
func (l *LLMPrompts) AddUserPrompt(new string) {
	l.User.WriteString(new + "\n")
}

type LLMResponse string

type LLMProvider interface {
	Generate(ctx context.Context, prompts *LLMPrompts, responseSchema json.RawMessage) (json.RawMessage, error)
}
