package agents

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/kakkky/hakoniwa/domain"
// 	"github.com/kakkky/hakoniwa/domain/mock"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"
// )

// // TestRuntime_Run_CancelStopsRun: Run は ctx cancel で nil を返して終了する
// func TestRuntime_Run_CancelStopsRun(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	llm := mock.NewMockLLMProvider(ctrl)
// 	sut := NewRuntime(llm)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	runErrCh := make(chan error, 1)
// 	go func() { runErrCh <- sut.Run(ctx) }()

// 	cancel()

// 	select {
// 	case err := <-runErrCh:
// 		if err != nil {
// 			t.Errorf("expected nil error from Run, got %v", err)
// 		}
// 	case <-time.After(time.Second):
// 		t.Fatal("Run did not return within 1s after ctx cancel")
// 	}
// }

// // TestCommandSubscriber_AddResidentAgent: AddResidentAgentCommand を受けると
// // reconcileResidentAgents が同期実行され、residentAgentsState.running に住人 agent が登録される
// func TestCommandSubscriber_AddResidentAgent(t *testing.T) {
// 	resident1, err := domain.NewResident("山田", 30, domain.Male, []domain.Trait{"優しい"})
// 	if err != nil {
// 		t.Fatalf("fixture: %v", err)
// 	}
// 	resident2, err := domain.NewResident("花子", 25, domain.Female, []domain.Trait{"明るい"})
// 	if err != nil {
// 		t.Fatalf("fixture: %v", err)
// 	}

// 	tests := []struct {
// 		name              string
// 		cmds              []domain.AgentCommand
// 		wantRunningIDs    []domain.ResidentID // residentAgentsState.running に期待する ID 集合
// 		wantBrokerRoutes  []domain.ResidentID // eventBroker.routes に期待する ID 集合
// 	}{
// 		{
// 			name: "1 件: running と routes に 1 件ずつ登録される",
// 			cmds: []domain.AgentCommand{
// 				domain.AddResidentAgentCommand{Resident: *resident1},
// 			},
// 			wantRunningIDs:   []domain.ResidentID{resident1.ID},
// 			wantBrokerRoutes: []domain.ResidentID{resident1.ID},
// 		},
// 		{
// 			name: "別住人 2 件: running と routes にそれぞれ 2 件登録される",
// 			cmds: []domain.AgentCommand{
// 				domain.AddResidentAgentCommand{Resident: *resident1},
// 				domain.AddResidentAgentCommand{Resident: *resident2},
// 			},
// 			wantRunningIDs:   []domain.ResidentID{resident1.ID, resident2.ID},
// 			wantBrokerRoutes: []domain.ResidentID{resident1.ID, resident2.ID},
// 		},
// 		{
// 			name: "同じ住人 2 件: running も routes も 1 件のまま (重複登録されない)",
// 			cmds: []domain.AgentCommand{
// 				domain.AddResidentAgentCommand{Resident: *resident1},
// 				domain.AddResidentAgentCommand{Resident: *resident1},
// 			},
// 			wantRunningIDs:   []domain.ResidentID{resident1.ID},
// 			wantBrokerRoutes: []domain.ResidentID{resident1.ID},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rt, sut, cancel := setupCommandSubscriberOnly(t)
// 			defer cancel()

// 			for _, cmd := range tt.cmds {
// 				sut.cmdInbox <- cmd
// 			}
// 			// reconcile は subscriber goroutine 内で同期実行されるが、テスト goroutine からは
// 			// 非同期に見える。雑だが 50ms 待てば十分処理される。
// 			time.Sleep(50 * time.Millisecond)

// 			rt.residentAgentsState.mu.Lock()
// 			gotRunning := make([]domain.ResidentID, 0, len(rt.residentAgentsState.running))
// 			for id := range rt.residentAgentsState.running {
// 				gotRunning = append(gotRunning, id)
// 			}
// 			rt.residentAgentsState.mu.Unlock()

// 			rt.eventBroker.routesMu.RLock()
// 			gotRoutes := make([]domain.ResidentID, 0, len(rt.eventBroker.routes))
// 			for id := range rt.eventBroker.routes {
// 				gotRoutes = append(gotRoutes, id)
// 			}
// 			rt.eventBroker.routesMu.RUnlock()

// 			assert.ElementsMatch(t, tt.wantRunningIDs, gotRunning, "running")
// 			assert.ElementsMatch(t, tt.wantBrokerRoutes, gotRoutes, "routes")
// 		})
// 	}
// }

// // TestCommandSubscriber_PublishEvent: PublishEventCommand を受けると Event が
// // eventBroker.inbox に転送される。eventBroker.run は起動しないので inbox は消費されない
// func TestCommandSubscriber_PublishEvent(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		cmd       domain.AgentCommand
// 		wantEvent domain.Event
// 	}{
// 		{
// 			name: "MessageEvent が同じ内容で eventBroker.inbox に転送される",
// 			cmd: domain.PublishEventCommand{
// 				Event: domain.MessageEvent{
// 					EventBase: domain.EventBase{
// 						EventTo:   domain.EventTo{ID: "resident-1", Name: "山田"},
// 						EventFrom: domain.EventFrom{ID: domain.BuildManagerActorID, Name: domain.BuoldManagerActorName},
// 					},
// 					Message: "hi",
// 				},
// 			},
// 			wantEvent: domain.MessageEvent{
// 				EventBase: domain.EventBase{
// 					EventTo:   domain.EventTo{ID: "resident-1", Name: "山田"},
// 					EventFrom: domain.EventFrom{ID: domain.BuildManagerActorID, Name: domain.BuoldManagerActorName},
// 				},
// 				Message: "hi",
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rt, sut, cancel := setupCommandSubscriberOnly(t)
// 			defer cancel()

// 			sut.cmdInbox <- tt.cmd

// 			select {
// 			case got := <-rt.eventBroker.inbox:
// 				assert.Equal(t, tt.wantEvent, got)
// 			case <-time.After(time.Second):
// 				t.Fatal("event was not forwarded to eventBroker.inbox within 1s")
// 			}
// 		})
// 	}
// }

// // TestRuntime_Run_EventBroker: 住人 1 人を起動した状態で宛先別の PublishEventCommand を流し、
// // 登録済み宛は届く / 未登録宛は捨てる、を観察する。
// //
// // なお配信の最終確認は residentAgent goroutine 内部に閉じるため state では観察できず、
// // やむなく LLM.Generate の呼び出しを proxy にしている。
// func TestRuntime_Run_EventBroker(t *testing.T) {
// 	resident, err := domain.NewResident("山田", 30, domain.Male, []domain.Trait{"優しい"})
// 	if err != nil {
// 		t.Fatalf("fixture: %v", err)
// 	}

// 	tests := []struct {
// 		name       string
// 		eventTo    domain.EventTo
// 		wantCalled bool
// 	}{
// 		{
// 			name:       "登録済み住人 ID 宛: 配信され LLM.Generate に到達する",
// 			eventTo:    domain.EventTo{ID: resident.ID, Name: resident.Name},
// 			wantCalled: true,
// 		},
// 		{
// 			name:       "未登録 ID 宛: 捨てられ LLM.Generate は呼ばれない",
// 			eventTo:    domain.EventTo{ID: "unknown"},
// 			wantCalled: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			llm := mock.NewMockLLMProvider(ctrl)
// 			generated := make(chan struct{}, 1)
// 			llm.EXPECT().
// 				Generate(gomock.Any(), gomock.Any()).
// 				DoAndReturn(func(_ context.Context, _ *domain.LLMPrompts) (domain.LLMResponse, error) {
// 					select {
// 					case generated <- struct{}{}:
// 					default:
// 					}
// 					return domain.LLMResponse(`{"action":"stay","payload":""}`), nil
// 				}).
// 				AnyTimes()

// 			sut := NewRuntime(llm)
// 			ctx, cancel := context.WithCancel(context.Background())
// 			defer cancel()
// 			go func() { _ = sut.Run(ctx) }()

// 			inbox := sut.CommandInbox()
// 			// 共通前提: 住人を 1 人起動しておく
// 			inbox <- domain.AddResidentAgentCommand{Resident: *resident}
// 			inbox <- domain.PublishEventCommand{
// 				Event: domain.MessageEvent{
// 					EventBase: domain.EventBase{
// 						EventTo:   tt.eventTo,
// 						EventFrom: domain.EventFrom{ID: domain.BuildManagerActorID, Name: domain.BuoldManagerActorName},
// 					},
// 					Message: "test",
// 				},
// 			}

// 			if tt.wantCalled {
// 				select {
// 				case <-generated:
// 				case <-time.After(2 * time.Second):
// 					t.Fatal("LLM.Generate was not called within 2s")
// 				}
// 			} else {
// 				select {
// 				case <-generated:
// 					t.Fatal("LLM.Generate should not have been called, but was")
// 				case <-time.After(300 * time.Millisecond):
// 				}
// 			}
// 		})
// 	}
// }

// func TestRuntime_CommandInbox(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		assert func(t *testing.T, sut *Runtime)
// 	}{
// 		{
// 			name: "非 nil の buffered チャネルを返す",
// 			assert: func(t *testing.T, sut *Runtime) {
// 				inbox := sut.CommandInbox()
// 				if inbox == nil {
// 					t.Fatal("CommandInbox returned nil channel")
// 				}
// 				select {
// 				case inbox <- domain.PublishEventCommand{}:
// 				default:
// 					t.Fatal("expected buffered channel, send was blocked")
// 				}
// 			},
// 		},
// 		{
// 			name: "複数回呼んでも同じチャネルを返す",
// 			assert: func(t *testing.T, sut *Runtime) {
// 				if sut.CommandInbox() != sut.CommandInbox() {
// 					t.Error("CommandInbox returned different channels across calls")
// 				}
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			llm := mock.NewMockLLMProvider(ctrl)
// 			sut := NewRuntime(llm)
// 			tt.assert(t, sut)
// 		})
// 	}
// }

// // --- helpers ---

// // setupCommandSubscriberOnly は Runtime を組み立て、commandSubscriber.run だけを起動する。
// // eventBroker.run / reconcileResidentAgentsLoop は起動しないため、
// // eventBroker.inbox は消費されず、外から観察できる。
// func setupCommandSubscriberOnly(t *testing.T) (rt *Runtime, sut *commandSubscriber, cancel context.CancelFunc) {
// 	t.Helper()
// 	ctrl := gomock.NewController(t)
// 	llm := mock.NewMockLLMProvider(ctrl)
// 	// 住人 agent が起動しても event が来なければ Generate は呼ばれないが、
// 	// 念のため呼ばれてもよい設定にする
// 	llm.EXPECT().
// 		Generate(gomock.Any(), gomock.Any()).
// 		Return(domain.LLMResponse(`{"action":"stay","payload":""}`), nil).
// 		AnyTimes()

// 	rt = NewRuntime(llm)
// 	sut = rt.commandSubscriber
// 	ctx, cancelFn := context.WithCancel(context.Background())
// 	go func() { _ = sut.run(ctx) }()
// 	return rt, sut, cancelFn
// }
