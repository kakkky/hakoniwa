package domain_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestCallLLM(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(llm *mock.MockLLMProvider)
		parse     func(domain.LLMResponse) (int, error)
		want      int
	}{
		{
			name: "Generate 成功 → parse 成功で結果が返る",
			setupMock: func(llm *mock.MockLLMProvider) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).Return(domain.LLMResponse("42"), nil)
			},
			parse: func(r domain.LLMResponse) (int, error) {
				return strconv.Atoi(string(r))
			},
			want: 42,
		},
		{
			name: "空文字応答でも parse が成功すれば値が返る",
			setupMock: func(llm *mock.MockLLMProvider) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).Return(domain.LLMResponse("0"), nil)
			},
			parse: func(r domain.LLMResponse) (int, error) {
				return strconv.Atoi(string(r))
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			llm := mock.NewMockLLMProvider(ctrl)
			tt.setupMock(llm)

			got, err := domain.CallLLM(context.Background(), llm, &domain.LLMPrompts{}, "schema", tt.parse)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got=%d want=%d", got, tt.want)
			}
		})
	}
}

func TestCallLLM_Error(t *testing.T) {
	generateErr := errors.New("llm down")
	parseErr := errors.New("parse failed")

	tests := []struct {
		name      string
		setupMock func(llm *mock.MockLLMProvider)
		parse     func(domain.LLMResponse) (int, error)
		wantErr   error
	}{
		{
			name: "Generate がエラーを返したらそのエラーが伝播する",
			setupMock: func(llm *mock.MockLLMProvider) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).Return(domain.LLMResponse(""), generateErr)
			},
			parse: func(_ domain.LLMResponse) (int, error) {
				// Generate がエラー時は parse が呼ばれない。呼ばれたらテスト失敗。
				return 0, errors.New("parse should not be called")
			},
			wantErr: generateErr,
		},
		{
			name: "parse がエラーを返したらそのエラーが伝播する",
			setupMock: func(llm *mock.MockLLMProvider) {
				llm.EXPECT().Generate(gomock.Any(), gomock.Any()).Return(domain.LLMResponse("not-a-number"), nil)
			},
			parse: func(_ domain.LLMResponse) (int, error) {
				return 0, parseErr
			},
			wantErr: parseErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			llm := mock.NewMockLLMProvider(ctrl)
			tt.setupMock(llm)

			got, err := domain.CallLLM(context.Background(), llm, &domain.LLMPrompts{}, "schema", tt.parse)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error: got=%v want=%v", err, tt.wantErr)
			}
			if got != 0 {
				t.Errorf("expected zero value, got %d", got)
			}
		})
	}
}
