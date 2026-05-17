package agents_test

import (
	"context"
	"testing"
	"time"

	"github.com/kakkky/hakoniwa/agents"
	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestRuntime_Run(t *testing.T) {
	tests := []struct {
		name   string
		act    func(t *testing.T, rt *agents.Runtime, cancel context.CancelFunc)
		assert func(t *testing.T, rt *agents.Runtime, runErrCh chan error)
	}{
		{
			name: "ctx を cancel すると Run が nil を返して終了する",
			act: func(_ *testing.T, _ *agents.Runtime, cancel context.CancelFunc) {
				cancel()
			},
			assert: func(t *testing.T, _ *agents.Runtime, runErrCh chan error) {
				select {
				case err := <-runErrCh:
					if err != nil {
						t.Errorf("expected nil error from Run, got %v", err)
					}
				case <-time.After(time.Second):
					t.Fatal("Run did not return within 1s after ctx cancel")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			llm := mock.NewMockLLMProvider(ctrl)
			rt := agents.NewRuntime(llm)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			runErrCh := make(chan error, 1)
			go func() { runErrCh <- rt.Run(ctx) }()

			tt.act(t, rt, cancel)
			tt.assert(t, rt, runErrCh)
		})
	}
}

func TestRuntime_CommandInbox(t *testing.T) {
	tests := []struct {
		name   string
		assert func(t *testing.T, rt *agents.Runtime)
	}{
		{
			name: "非 nil の buffered チャネルを返す",
			assert: func(t *testing.T, rt *agents.Runtime) {
				inbox := rt.CommandInbox()
				if inbox == nil {
					t.Fatal("CommandInbox returned nil channel")
				}
				select {
				case inbox <- domain.PublishEventCommand{}:
				default:
					t.Fatal("expected buffered channel, send was blocked")
				}
			},
		},
		{
			name: "複数回呼んでも同じチャネルを返す",
			assert: func(t *testing.T, rt *agents.Runtime) {
				if rt.CommandInbox() != rt.CommandInbox() {
					t.Error("CommandInbox returned different channels across calls")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			llm := mock.NewMockLLMProvider(ctrl)
			rt := agents.NewRuntime(llm)
			tt.assert(t, rt)
		})
	}
}
