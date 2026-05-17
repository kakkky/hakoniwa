package domain_test

import (
	"testing"

	"github.com/kakkky/hakoniwa/domain"
)

func TestNewResident(t *testing.T) {
	tests := []struct {
		name       string
		inputName  string
		inputAge   int
		inputGen   domain.Gender
		inputTraits []domain.Trait
	}{
		{
			name:       "全フィールドを正しく渡せば作成できる",
			inputName:  "山田",
			inputAge:   30,
			inputGen:   domain.Male,
			inputTraits: []domain.Trait{"優しい"},
		},
		{
			name:       "age が 0 でも作成できる",
			inputName:  "幼児",
			inputAge:   0,
			inputGen:   domain.Female,
			inputTraits: []domain.Trait{"元気"},
		},
		{
			name:       "gender が境界値 Unspecified でも作成できる",
			inputName:  "名無し",
			inputAge:   20,
			inputGen:   domain.Unspecified,
			inputTraits: []domain.Trait{"無口"},
		},
		{
			name:       "gender が境界値 Female でも作成できる",
			inputName:  "花子",
			inputAge:   25,
			inputGen:   domain.Female,
			inputTraits: []domain.Trait{"明るい", "几帳面"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewResident(tt.inputName, tt.inputAge, tt.inputGen, tt.inputTraits)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("expected non-nil resident")
			}
			if string(got.Name) != tt.inputName {
				t.Errorf("Name: got=%q want=%q", got.Name, tt.inputName)
			}
			if got.Age != tt.inputAge {
				t.Errorf("Age: got=%d want=%d", got.Age, tt.inputAge)
			}
			if got.Gender != tt.inputGen {
				t.Errorf("Gender: got=%d want=%d", got.Gender, tt.inputGen)
			}
			if len(got.Traits) != len(tt.inputTraits) {
				t.Errorf("Traits length: got=%d want=%d", len(got.Traits), len(tt.inputTraits))
			}
			if got.ID == "" {
				t.Error("ID should be non-empty")
			}
		})
	}
}

func TestNewResident_Error(t *testing.T) {
	tests := []struct {
		name       string
		inputName  string
		inputAge   int
		inputGen   domain.Gender
		inputTraits []domain.Trait
		wantMsg    string
	}{
		{
			name:       "name が空文字ならエラー",
			inputName:  "",
			inputAge:   30,
			inputGen:   domain.Male,
			inputTraits: []domain.Trait{"優しい"},
			wantMsg:    "name is required",
		},
		{
			name:       "age が負ならエラー",
			inputName:  "山田",
			inputAge:   -1,
			inputGen:   domain.Male,
			inputTraits: []domain.Trait{"優しい"},
			wantMsg:    "age must be non-negative",
		},
		{
			name:       "gender が Unspecified より小さいとエラー",
			inputName:  "山田",
			inputAge:   30,
			inputGen:   domain.Gender(-1),
			inputTraits: []domain.Trait{"優しい"},
			wantMsg:    "invalid gender",
		},
		{
			name:       "gender が Female を超えるとエラー",
			inputName:  "山田",
			inputAge:   30,
			inputGen:   domain.Gender(3),
			inputTraits: []domain.Trait{"優しい"},
			wantMsg:    "invalid gender",
		},
		{
			name:       "traits が空ならエラー",
			inputName:  "山田",
			inputAge:   30,
			inputGen:   domain.Male,
			inputTraits: []domain.Trait{},
			wantMsg:    "at least one trait is required",
		},
		{
			name:       "traits が nil ならエラー",
			inputName:  "山田",
			inputAge:   30,
			inputGen:   domain.Male,
			inputTraits: nil,
			wantMsg:    "at least one trait is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewResident(tt.inputName, tt.inputAge, tt.inputGen, tt.inputTraits)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if got != nil {
				t.Errorf("expected nil resident, got %+v", got)
			}
			if err.Error() != tt.wantMsg {
				t.Errorf("error message: got=%q want=%q", err.Error(), tt.wantMsg)
			}
		})
	}
}
