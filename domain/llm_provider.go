package domain

//go:generate mockgen -source=llm_provider.go -destination=../testing/mock/mock_llm_provider.go -package=mock

import (
	"context"
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
	Generate(ctx context.Context, prompts *LLMPrompts) (LLMResponse, error)
}

func CallLLM[RES any](
	ctx context.Context,
	p LLMProvider,
	prompts *LLMPrompts,
	schema string,
	parse func(LLMResponse) (RES, error),
) (RES, error) {
	var zero RES
	raw, err := p.Generate(ctx, prompts)
	if err != nil {
		return zero, err
	}
	return parse(raw)
}
