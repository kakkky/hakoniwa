package llm

import (
	"context"
	"encoding/json"

	"github.com/kakkky/hakoniwa/config"
	"github.com/kakkky/hakoniwa/domain"
	"google.golang.org/genai"
)

const model = ""

type LLMGeminiProvider struct {
	client *genai.Client
}

func NewLLMGeminiProvider(ctx context.Context, cfg config.Config) (*LLMGeminiProvider, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: cfg.GeminiAPIKey,
	})
	if err != nil {
		return nil, err
	}
	return &LLMGeminiProvider{
		client: client,
	}, nil
}

func (p *LLMGeminiProvider) Generate(ctx context.Context, prompts *domain.LLMPrompts, responseSchema *json.RawMessage) (*json.RawMessage, error) {
	p.client.Models.GenerateContent(ctx, model, nil, nil)
	return nil, nil
}
