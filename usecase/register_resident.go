package usecase

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/kakkky/hakoniwa/domain"
	llmresponse "github.com/kakkky/hakoniwa/schema/llm_response"
)

type RegisterResident struct {
	repository            domain.ResidentRepository
	llmProvider           domain.LLMProvider
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

type registerResidentExtractTraitsResponse struct {
	Traits []domain.Trait `json:"traits"`
}

func (r *RegisterResident) Exec(ctx context.Context, name string, age int, gender domain.Gender, personalityDescription string) error {
	userPromptTemplate :=
		`以下の文章はある人の性格を表した文章です。その文章をいくつかの特徴(traits)に落とし込んでください
		文章：%s
		`

	userPrompt := fmt.Sprintf(userPromptTemplate, personalityDescription)
	var llmPrompt domain.LLMPrompts
	llmPrompt.AddUserPrompt(userPrompt)

	rawResp, err := r.llmProvider.Generate(ctx, &llmPrompt, llmresponse.RegisterResidentExtractTraits)
	if err != nil {
		return fmt.Errorf("failed to generate traits: %w", err)
	}
	var registerResidentExtractTraits registerResidentExtractTraitsResponse
	if err := json.Unmarshal(rawResp, &registerResidentExtractTraits); err != nil {
		return nil
	}

	traits := registerResidentExtractTraits.Traits

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
