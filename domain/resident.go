// Package domain is domain layer of the app
package domain

import (
	"errors"

	"github.com/oklog/ulid/v2"
)

type Resident struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Gender Gender  `json:"gender"`
	Traits []Trait `json:"traints"`
}

type Gender int

const (
	Unspecified Gender = iota
	Male        Gender = iota
	Female      Gender = iota
)

type Trait string

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
		ID:     id,
		Name:   name,
		Age:    age,
		Gender: gender,
		Traits: traints,
	}, nil
}
