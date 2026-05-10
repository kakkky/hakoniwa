package domain

import "context"

// Reposutory Interfaces
type ResidentRepository interface {
	Save(resident *Resident) error
	// FindByID(id string) (*Resident, error)
	// DeleteByID(id string) error
}

// LLM Provider Interface
type Prompts struct {
	System string
	User   string
}

type LLMResponse string

type LLMProvider interface {
	Generate(ctx context.Context, prompts Prompts) (LLMResponse, error)
}
