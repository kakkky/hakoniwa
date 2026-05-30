package llm

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/kakkky/hakoniwa/config"
	"github.com/kakkky/hakoniwa/domain"
	"google.golang.org/genai"
)

const model = "gemini-2.5-flash-lite"

type LLMGeminiProvider struct {
	client *genai.Client
}

func NewLLMGeminiProvider(cfg *config.Config) (*LLMGeminiProvider, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: cfg.GeminiAPIKey,
	})
	if err != nil {
		return nil, err
	}
	return &LLMGeminiProvider{
		client: client,
	}, nil
}

func (p *LLMGeminiProvider) Generate(ctx context.Context, prompts domain.LLMPrompts, responseSchema json.RawMessage) (json.RawMessage, error) {
	if responseSchema == nil {
		return nil, errors.New("schema is required")
	}
	var systemInstruction *genai.Content
	if prompts.System != "" {
		systemInstruction = &genai.Content{
			Parts: []*genai.Part{
				{Text: prompts.System},
			},
		}
	}
	gcc := &genai.GenerateContentConfig{
		SystemInstruction:  systemInstruction,
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: responseSchema,
	}
	contents := []*genai.Content{
		genai.NewContentFromText(prompts.User, genai.RoleUser),
	}
	time.Sleep(1 * time.Second)
	resp, err := p.client.Models.GenerateContent(ctx, model, contents, gcc)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(resp.Text()), nil
}
