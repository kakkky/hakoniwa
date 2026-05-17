package agents

import (
	"context"
	"sync"

	"github.com/kakkky/hakoniwa/domain"
)

type Runtime struct {
	gameMasterAgent     *gameMasterAgent
	residentAgentsState residentAgentsState

	commandSubscriber *commandSubscriber
	eventBroker       *eventBroker

	llmProvider domain.LLMProvider

	reconcileSignal reconcileSignal
}

type residentAgentsState struct {
	mu      sync.Mutex
	desired []*residentAgent
	running map[domain.ResidentID]context.CancelFunc
}

func NewRuntime(llmProvider domain.LLMProvider) *Runtime {
	r := &Runtime{
		eventBroker:     newEventBroker(),
		llmProvider:     llmProvider,
		reconcileSignal: make(reconcileSignal, 1),
		residentAgentsState: residentAgentsState{
			running: make(map[domain.ResidentID]context.CancelFunc),
		},
	}
	r.commandSubscriber = newCommandSubscriber(r)
	return r
}

func (r *Runtime) Run(ctx context.Context) error {
	go r.eventBroker.run(ctx)
	go r.commandSubscriber.run(ctx)

	r.reconcileResidentAgentsLoop(ctx)

	return nil
}

// CommandInbox は外部から command を流し込むためのチャネルを返す。
// AgentCommander の実装にこのチャネルを渡して利用する。
func (r *Runtime) CommandInbox() domain.AgentCommandInbox {
	return r.commandSubscriber.cmdInbox
}

type reconcileSignal chan struct{}

func (r *Runtime) reconcileResidentAgentsLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-r.reconcileSignal:
			r.reconcileResidentAgents(ctx)
		}
	}
}

func (r *Runtime) reconcileResidentAgents(ctx context.Context) {
	r.residentAgentsState.mu.Lock()
	defer r.residentAgentsState.mu.Unlock()
	for _, desired := range r.residentAgentsState.desired {
		if _, alreadyRunning := r.residentAgentsState.running[desired.resident.ID]; alreadyRunning {
			continue
		}
		agentCtx, cancelFn := context.WithCancel(ctx)

		sendEvent := func(event domain.Event) { r.eventBroker.inbox <- event }

		residentAgent := newResidentAgent(newAgentBase(sendEvent, r.llmProvider), desired.resident)
		r.residentAgentsState.running[residentAgent.resident.ID] = cancelFn

		r.eventBroker.registerRoutes(residentAgent.resident.ID, residentAgent.inbox)

		go residentAgent.run(agentCtx)
	}
}

func (r *Runtime) addResidentAgent(resident *domain.Resident) {
	r.residentAgentsState.mu.Lock()
	r.residentAgentsState.desired = append(r.residentAgentsState.desired, &residentAgent{
		resident: resident,
	})
	r.residentAgentsState.mu.Unlock()
	select {
	case r.reconcileSignal <- struct{}{}:
	default:
	}
}
