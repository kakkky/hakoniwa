package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/domain/mock"
	"github.com/kakkky/hakoniwa/usecase"
	"go.uber.org/mock/gomock"
)

func TestRegisterResident_Exec(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
		inputAge  int
		inputGen  domain.Gender
		setupMock func(llm *mock.MockLLMProvider, repo *mock.MockResidentRepository, cmder *mock.MockAgentCommander, inputName string, inputAge int, inputGen domain.Gender)
	}{
		{
			name:      "Generate→Save→PublishCommand が全て成功すれば nil",
			inputName: "山田",
			inputAge:  30,
			inputGen:  domain.Male,
			setupMock: func(llm *mock.MockLLMProvider, repo *mock.MockResidentRepository, cmder *mock.MockAgentCommander, inputName string, inputAge int, inputGen domain.Gender) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return(domain.LLMResponse(`{"traits":["優しい","几帳面"]}`), nil)
				repo.EXPECT().Save(gomock.Cond(func(r *domain.Resident) bool {
					return r != nil &&
						string(r.Name) == inputName &&
						r.Age == inputAge &&
						r.Gender == inputGen &&
						len(r.Traits) > 0
				})).Return(nil)
				cmder.EXPECT().PublishCommand(
					gomock.Any(),
					gomock.AssignableToTypeOf(domain.AddResidentAgentCommand{}),
				).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			llm := mock.NewMockLLMProvider(ctrl)
			repo := mock.NewMockResidentRepository(ctrl)
			cmder := mock.NewMockAgentCommander(ctrl)
			tt.setupMock(llm, repo, cmder, tt.inputName, tt.inputAge, tt.inputGen)

			uc := usecase.NewRegisterResident(repo, llm, cmder)
			if err := uc.Exec(context.Background(), tt.inputName, tt.inputAge, tt.inputGen, "穏やかで几帳面"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRegisterResident_Exec_Error(t *testing.T) {
	generateErr := errors.New("llm down")
	saveErr := errors.New("save failed")
	publishErr := errors.New("publish failed")

	tests := []struct {
		name               string
		inputName          string
		inputAge           int
		inputGen           domain.Gender
		setupMock          func(llm *mock.MockLLMProvider, repo *mock.MockResidentRepository, cmder *mock.MockAgentCommander)
		wantErrIs          error
		wantErrMsgContains string
	}{
		{
			name:      "Generate がエラーを返すと wrap して返す",
			inputName: "山田",
			inputAge:  30,
			inputGen:  domain.Male,
			setupMock: func(llm *mock.MockLLMProvider, _ *mock.MockResidentRepository, _ *mock.MockAgentCommander) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).Return(domain.LLMResponse(""), generateErr)
			},
			wantErrIs:          generateErr,
			wantErrMsgContains: "failed to generate traits",
		},
		{
			name:      "Generate が invalid JSON を返すと parse エラーが wrap される",
			inputName: "山田",
			inputAge:  30,
			inputGen:  domain.Male,
			setupMock: func(llm *mock.MockLLMProvider, _ *mock.MockResidentRepository, _ *mock.MockAgentCommander) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).Return(domain.LLMResponse("not json"), nil)
			},
			wantErrMsgContains: "failed to generate traits",
		},
		{
			name:      "name が空だと NewResident validation エラーが wrap される",
			inputName: "",
			inputAge:  30,
			inputGen:  domain.Male,
			setupMock: func(llm *mock.MockLLMProvider, _ *mock.MockResidentRepository, _ *mock.MockAgentCommander) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).Return(domain.LLMResponse(`{"traits":["優しい"]}`), nil)
			},
			wantErrMsgContains: "failed to create resident",
		},
		{
			name:      "Save がエラーを返すと wrap して返す",
			inputName: "山田",
			inputAge:  30,
			inputGen:  domain.Male,
			setupMock: func(llm *mock.MockLLMProvider, repo *mock.MockResidentRepository, _ *mock.MockAgentCommander) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).Return(domain.LLMResponse(`{"traits":["優しい"]}`), nil)
				repo.EXPECT().Save(gomock.Any()).Return(saveErr)
			},
			wantErrIs:          saveErr,
			wantErrMsgContains: "failed to save resident",
		},
		{
			name:      "PublishCommand がエラーを返すと wrap して返す",
			inputName: "山田",
			inputAge:  30,
			inputGen:  domain.Male,
			setupMock: func(llm *mock.MockLLMProvider, repo *mock.MockResidentRepository, cmder *mock.MockAgentCommander) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).Return(domain.LLMResponse(`{"traits":["優しい"]}`), nil)
				repo.EXPECT().Save(gomock.Any()).Return(nil)
				cmder.EXPECT().PublishCommand(gomock.Any(), gomock.Any()).Return(publishErr)
			},
			wantErrIs:          publishErr,
			wantErrMsgContains: "failed to publish AddResidentAgentCommand",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			llm := mock.NewMockLLMProvider(ctrl)
			repo := mock.NewMockResidentRepository(ctrl)
			cmder := mock.NewMockAgentCommander(ctrl)
			tt.setupMock(llm, repo, cmder)

			uc := usecase.NewRegisterResident(repo, llm, cmder)
			err := uc.Exec(context.Background(), tt.inputName, tt.inputAge, tt.inputGen, "personality")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("errors.Is mismatch: got=%v want=%v", err, tt.wantErrIs)
			}
			if tt.wantErrMsgContains != "" && !strings.Contains(err.Error(), tt.wantErrMsgContains) {
				t.Errorf("error message %q should contain %q", err.Error(), tt.wantErrMsgContains)
			}
		})
	}
}
