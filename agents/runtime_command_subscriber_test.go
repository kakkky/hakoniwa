package agents

import (
	"context"
	"testing"
	"time"

	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestCommandSubscriber_Run(t *testing.T) {
	residentFixture, err := domain.NewResident("山田", 30, domain.Male, []domain.Trait{"優しい"})
	if err != nil {
		t.Fatalf("fixture setup: %v", err)
	}

	tests := []struct {
		name        string
		sendCmd     domain.AgentCommand // cmdInbox に送るコマンド (nil なら送らない)
		cancelInAct bool                // act で ctx cancel するか
		assert      func(t *testing.T, rt *Runtime, runErrCh chan error)
	}{
		{
			name:        "ctx を cancel すると run が nil を返して終了する",
			cancelInAct: true,
			assert: func(t *testing.T, _ *Runtime, runErrCh chan error) {
				select {
				case err := <-runErrCh:
					if err != nil {
						t.Errorf("expected nil error, got %v", err)
					}
				case <-time.After(time.Second):
					t.Fatal("run did not return within 1s after ctx cancel")
				}
			},
		},
		{
			name:    "AddResidentAgentCommand を受信すると runtime.desired に追加され reconcileSignal が送られる",
			sendCmd: domain.AddResidentAgentCommand{Resident: *residentFixture},
			assert: func(t *testing.T, rt *Runtime, _ chan error) {
				// reconcileSignal の受信を append 完了の同期点として使う
				select {
				case <-rt.reconcileSignal:
				case <-time.After(time.Second):
					t.Fatal("reconcileSignal was not sent within 1s")
				}

				rt.residentAgentsState.mu.Lock()
				gotLen := len(rt.residentAgentsState.desired)
				var gotID domain.ResidentID
				if gotLen == 1 {
					gotID = rt.residentAgentsState.desired[0].resident.ID
				}
				rt.residentAgentsState.mu.Unlock()

				if gotLen != 1 {
					t.Fatalf("desired length: got=%d want=1", gotLen)
				}
				if gotID != residentFixture.ID {
					t.Errorf("desired[0].resident.ID: got=%q want=%q", gotID, residentFixture.ID)
				}
			},
		},
		{
			name: "PublishEventCommand を受信すると eventBroker.inbox に転送される",
			sendCmd: domain.PublishEventCommand{
				Event: domain.MessageEvent{
					EventBase: domain.EventBase{
						EventTo:   domain.EventTo{ID: "x"},
						EventFrom: domain.EventFrom{ID: domain.BuildManagerActorID},
					},
					Message: "hi",
				},
			},
			assert: func(t *testing.T, rt *Runtime, _ chan error) {
				select {
				case got := <-rt.eventBroker.inbox:
					me, ok := got.(domain.MessageEvent)
					if !ok {
						t.Fatalf("expected MessageEvent, got %T", got)
					}
					if me.Message != "hi" {
						t.Errorf("Message: got=%q want=%q", me.Message, "hi")
					}
				case <-time.After(time.Second):
					t.Fatal("event was not forwarded to eventBroker.inbox within 1s")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			llm := mock.NewMockLLMProvider(ctrl)
			rt := NewRuntime(llm)
			sub := rt.commandSubscriber

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			runErrCh := make(chan error, 1)
			go func() { runErrCh <- sub.run(ctx) }()

			if tt.sendCmd != nil {
				sub.cmdInbox <- tt.sendCmd
			}
			if tt.cancelInAct {
				cancel()
			}

			tt.assert(t, rt, runErrCh)
		})
	}
}
