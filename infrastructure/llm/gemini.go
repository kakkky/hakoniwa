package llm

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type LLMGeminiProvider struct{}

func NewLLMGeminiProvider() domain.LLMProvider {
	return &LLMGeminiProvider{}
}

func (p *LLMGeminiProvider) Generate(ctx context.Context, prompts *domain.LLMPrompts, schema string) (domain.LLMResponse, error) {
	panic("not implemented")
}
