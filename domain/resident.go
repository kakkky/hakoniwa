// Package domain is domain layer of the app
package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

type ResidentID string

type Resident struct {
	ID       ResidentID `json:"id"`
	Name     string     `json:"name"`
	Age      int        `json:"age"`
	Gender   Gender     `json:"gender"`
	Traits   []Trait    `json:"traints"`
	Mood     Mood       `json:"mood"`
	Memories []Memory   `json:"memories"`
}

type Gender int

const (
	Unspecified Gender = iota
	Male        Gender = iota
	Female      Gender = iota
)

type Trait string

type Mood string

type Memory struct {
	Content   string    `json:"content"`
	OccuredAt time.Time `json:"occured_at"`
}

func (m Memory) String() string {
	return fmt.Sprintf("%s: %s", m.OccuredAt, m.Content)
}

func NewResident(name string, age int, gender Gender, traints []Trait) (*Resident, error) {
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
	if len(traints) == 0 {
		return nil, errors.New("at least one trait is required")
	}

	id := ulid.Make().String()
	return &Resident{
		ID:     ResidentID(id),
		Name:   name,
		Age:    age,
		Gender: gender,
		Traits: traints,
	}, nil
}

func (r *Resident) UpdateMood(mood Mood) {
	r.Mood = mood
}

func (r *Resident) AddMemory(memories ...Memory) {
	r.Memories = append(r.Memories, memories...)
}
