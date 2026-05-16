package domain

import (
	"context"
	"strings"
)

// LLM Provider Interface
type LLMPrompts struct {
	System strings.Builder
	User   strings.Builder
}

func (l *LLMPrompts) AddSystemPromt(new string) {
	l.System.WriteString(new + "\n")
}
func (l *LLMPrompts) AddUserPromt(new string) {
	l.User.WriteString(new + "\n")
}

type LLMResponse string

type LLMProvider interface {
	Generate(ctx context.Context, prompts *LLMPrompts, schema string) (LLMResponse, error)
}

func CallLLM[RES any](
	ctx context.Context,
	p LLMProvider,
	prompts *LLMPrompts,
	schema string,
	parse func(LLMResponse) (RES, error),
) (RES, error) {
	var zero RES
	raw, err := p.Generate(ctx, prompts, schema)
	if err != nil {
		return zero, err
	}
	return parse(raw)
}
