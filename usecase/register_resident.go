package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/kakkky/hakoniwa/domain"
)

type RegisterResident struct {
	repository  domain.ResidentRepository
	llmProvider domain.LLMProvider
}

func NewRegisterResident(repository domain.ResidentRepository, llmProvider domain.LLMProvider) *RegisterResident {
	return &RegisterResident{
		repository:  repository,
		llmProvider: llmProvider,
	}
}

func (r *RegisterResident) Exec(ctx context.Context, name string, age int, gender domain.Gender, personalityDescription string) error {
	systemPrompt := `


	`
	userPromptTemplate := `

	`
	userPrompt := fmt.Sprintf(userPromptTemplate, personalityDescription)
	var llmPrompt domain.LLMPrompts
	llmPrompt.AddSystemPrompt(systemPrompt)
	llmPrompt.AddUserPrompt(userPrompt)

	traits, err := domain.CallLLM(ctx, r.llmProvider, &llmPrompt, "", parseTraits)
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

	return nil
}

func parseTraits(raw domain.LLMResponse) ([]domain.Trait, error) {
	rawTraitSlice := strings.Split(string(raw), ",")
	traits := make([]domain.Trait, len(rawTraitSlice))
	for i, rawTrait := range rawTraitSlice {
		traits[i] = domain.Trait(strings.TrimSpace(rawTrait))
	}
	return traits, nil
}
