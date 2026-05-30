// Package domain is domain layer of the app
package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

type ResidentID string
type ResidentName string

type Resident struct {
	ID       ResidentID   `json:"id"`
	Name     ResidentName `json:"name"`
	Age      int          `json:"age"`
	Gender   Gender       `json:"gender"`
	Traits   Traits       `json:"traits"`
	Mood     Mood         `json:"mood"`
	Memories Memories     `json:"memories"`
}

type Gender int

func (g Gender) String() string {
	switch g {
	case Male:
		return "男性"
	case Female:
		return "女性"
	default:
		return "不明"
	}
}

const (
	Unspecified Gender = iota
	Male        Gender = iota
	Female      Gender = iota
)

type Trait string
type Traits []Trait

func (ts Traits) String() string {
	traits := make([]string, 0, len(ts))
	for _, t := range ts {
		traits = append(traits, string(t))
	}
	return strings.Join(traits, ",")
}

type Mood string

type Memory struct {
	Content   string    `json:"content"`
	OccuredAt time.Time `json:"occured_at"`
}

type Memories []Memory

func (ms Memories) String() string {
	memories := make([]string, len(ms))
	for _, m := range ms {
		memory := fmt.Sprintf("【%s】%s", m.OccuredAt, m.Content)
		memories = append(memories, memory)
	}
	return strings.Join(memories, ",")
}

func NewResident(name string, age int, gender Gender, traits []Trait) (*Resident, error) {
	if name == "" {
		// のちに適当なdomain errorを定義して返すべき
		return nil, errors.New("name is required")
	}
	if age < 0 {
		return nil, errors.New("age must be non-negative")
	}
	if gender < Unspecified || gender > Female {
		return nil, errors.New("invalid gender")
	}
	if len(traits) == 0 {
		return nil, errors.New("at least one trait is required")
	}

	id := ulid.Make().String()
	return &Resident{
		ID:     ResidentID(id),
		Name:   ResidentName(name),
		Age:    age,
		Gender: gender,
		Traits: traits,
	}, nil
}

func (r *Resident) UpdateMood(mood Mood) {
	r.Mood = mood
}

func (r *Resident) AddMemory(memories ...Memory) {
	r.Memories = append(r.Memories, memories...)
}
