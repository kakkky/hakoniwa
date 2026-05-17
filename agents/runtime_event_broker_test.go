package agents

import (
	"context"
	"testing"
	"time"

	"github.com/kakkky/hakoniwa/domain"
)

func TestEventBroker_Run(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(b *eventBroker) agentEventInbox // 登録した route inbox を返す (不要なら nil)
		act    func(b *eventBroker, cancel context.CancelFunc)
		assert func(t *testing.T, registered agentEventInbox, runErrCh chan error)
	}{
		{
			name: "ctx を cancel すると run が nil を返して終了する",
			act: func(_ *eventBroker, cancel context.CancelFunc) {
				cancel()
			},
			assert: func(t *testing.T, _ agentEventInbox, runErrCh chan error) {
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
			name: "inbox に流した event は registerRoutes で登録された inbox に届く",
			setup: func(b *eventBroker) agentEventInbox {
				inbox := make(agentEventInbox, 1)
				b.registerRoutes("resident-1", inbox)
				return inbox
			},
			act: func(b *eventBroker, _ context.CancelFunc) {
				b.inbox <- domain.MessageEvent{
					EventBase: domain.EventBase{
						EventTo:   domain.EventTo{ID: "resident-1", Name: "山田"},
						EventFrom: domain.EventFrom{ID: domain.BuildManagerActorID, Name: domain.BuoldManagerActorName},
					},
					Message: "hello",
				}
			},
			assert: func(t *testing.T, registered agentEventInbox, _ chan error) {
				select {
				case got := <-registered:
					me, ok := got.(domain.MessageEvent)
					if !ok {
						t.Fatalf("expected MessageEvent, got %T", got)
					}
					if me.Message != "hello" {
						t.Errorf("Message: got=%q want=%q", me.Message, "hello")
					}
				case <-time.After(time.Second):
					t.Fatal("event was not delivered within 1s")
				}
			},
		},
		{
			name: "未登録 ID 宛の event は捨てられ、登録済み宛 event はその後も届く",
			setup: func(b *eventBroker) agentEventInbox {
				inbox := make(agentEventInbox, 1)
				b.registerRoutes("known", inbox)
				return inbox
			},
			act: func(b *eventBroker, _ context.CancelFunc) {
				// 1 件目: 未登録 ID 宛 → 捨てられる
				b.inbox <- domain.MessageEvent{
					EventBase: domain.EventBase{EventTo: domain.EventTo{ID: "unknown"}},
					Message:   "drop me",
				}
				// 2 件目: 登録済み ID 宛 → 届く
				b.inbox <- domain.MessageEvent{
					EventBase: domain.EventBase{EventTo: domain.EventTo{ID: "known"}},
					Message:   "deliver me",
				}
			},
			assert: func(t *testing.T, registered agentEventInbox, _ chan error) {
				select {
				case got := <-registered:
					if me, ok := got.(domain.MessageEvent); !ok || me.Message != "deliver me" {
						t.Errorf("unexpected event: %+v", got)
					}
				case <-time.After(time.Second):
					t.Fatal("registered event was not delivered (broker may have deadlocked)")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newEventBroker()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			runErrCh := make(chan error, 1)
			go func() { runErrCh <- b.run(ctx) }()

			var registered agentEventInbox
			if tt.setup != nil {
				registered = tt.setup(b)
			}
			tt.act(b, cancel)
			tt.assert(t, registered, runErrCh)
		})
	}
}
