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
	userPromtTemplate := `

	`
	userPrompt := fmt.Sprintf(userPromtTemplate, personalityDescription)

	response, err := r.llmProvider.Generate(ctx, domain.Prompts{System: systemPrompt, User: userPrompt})
	if err != nil {
		return fmt.Errorf("failed to generate traits: %w", err)
	}
	traits := parseTraits(string(response))

	resident, err := domain.NewResident(name, age, gender, traits)
	if err != nil {
		return fmt.Errorf("failed to create resident: %w", err)
	}

	if err := r.repository.Save(resident); err != nil {
		return fmt.Errorf("failed to save resident: %w", err)
	}

	return nil
}

func parseTraits(rawTraits string) []domain.Trait {
	rawTraitSlice := strings.Split(rawTraits, ",")
	traits := make([]domain.Trait, len(rawTraitSlice))
	for i, rawTrait := range rawTraitSlice {
		traits[i] = domain.Trait(strings.TrimSpace(rawTrait))
	}
	return traits
}
