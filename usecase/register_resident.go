package usecase

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/kakkky/hakoniwa/domain"
)

type RegisterResident struct {
	repository     domain.ResidentRepository
	llmProvider    domain.LLMProvider
	agentCommandPublisher domain.AgentCommandPublisher
}

func NewRegisterResident(
	repository domain.ResidentRepository,
	llmProvider domain.LLMProvider,
	agentCommandPublisher domain.AgentCommandPublisher,
) *RegisterResident {
	return &RegisterResident{
		repository:            repository,
		llmProvider:           llmProvider,
		agentCommandPublisher: agentCommandPublisher,
	}
}

func (r *RegisterResident) Exec(ctx context.Context, name string, age int, gender domain.Gender, personalityDescription string) error {
	systemPromptTemplate := `レスポンスは以下のJSONスキーマの文字列で行うようにしてください。：%s`
	systemPrompt := fmt.Sprintf(systemPromptTemplate, registerResidentLLMResponseSchema)
	userPromptTemplate := `

	`
	userPrompt := fmt.Sprintf(userPromptTemplate, personalityDescription)
	var llmPrompt domain.LLMPrompts
	llmPrompt.AddSystemPrompt(systemPrompt)
	llmPrompt.AddUserPrompt(userPrompt)

	traits, err := domain.CallLLM(ctx, r.llmProvider, &llmPrompt, registerResidentLLMResponseSchema, parseRegisterResidentLLMResponse)
	if err != nil {
		return fmt.Errorf("failed to generate traits: %w", err)
	}

	resident, err := domain.NewResident(name, age, gender, traits)
	if err != nil {
		return fmt.Errorf("failed to create resident: %w", err)
	}

	if err := r.repository.Save(resident); err != nil {
		return fmt.Errorf("failed to save resident: %w", err)
	}

	if err := r.agentCommandPublisher.PublishCommand(ctx, domain.AddResidentAgentCommand{Resident: *resident}); err != nil {
		return fmt.Errorf("failed to publish AddResidentAgentCommand: %w", err)
	}

	return nil
}

//go:embed schema/register_resident_llm_response_schema.json
var registerResidentLLMResponseSchema string

type registerResidentLLMResponse struct {
	Traits []string `json:"traits"`
}

func parseRegisterResidentLLMResponse(raw domain.LLMResponse) ([]domain.Trait, error) {
	var res registerResidentLLMResponse
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return nil, err
	}
	traits := make([]domain.Trait, len(res.Traits))
	for i, t := range res.Traits {
		traits[i] = domain.Trait(t)
	}
	return traits, nil
}
