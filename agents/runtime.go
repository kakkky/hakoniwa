package agents

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type Runtime struct {
	gameMasterAgent     *gameMasterAgent
	residentAgentsState residentAgentsState

	commandSubscriber commandSubscriber
	eventBroker       eventBroker

	llmProvider domain.LLMProvider

	reconcileSignal reconcileSignal
}

type residentAgentsState struct {
	desired []*residentAgent
	running map[id]context.CancelFunc
}

func NewRuntime(
	commandSubscriber commandSubscriber,
	eventBroker eventBroker,
) *Runtime {
	return &Runtime{
		commandSubscriber: commandSubscriber,
		eventBroker:       eventBroker,
	}
}

func (r *Runtime) Run(ctx context.Context) error {
	go r.eventBroker.run(ctx)
	go r.commandSubscriber.run()

	r.reconcileResidentAgentsLoop(ctx)

	return nil
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
	for _, desired := range r.residentAgentsState.desired {
		if _, alreadyRunning := r.residentAgentsState.running[desired.id]; alreadyRunning {
			continue
		}
		agentCtx, cancelFn := context.WithCancel(ctx)

		sendEvent := func(event agentEvent) { r.eventBroker.inbox <- event }

		residentAgent := newResidentAgent(newAgentBase(sendEvent, r.llmProvider), desired.resident)
		r.residentAgentsState.running[residentAgent.id] = cancelFn

		r.eventBroker.registerRoutes(residentAgent.id, residentAgent.inbox)

		go residentAgent.run(agentCtx)
	}
}

func (r *Runtime) addResidentAgent(resident *domain.Resident) {

}
