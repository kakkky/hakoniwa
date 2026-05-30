package agents

import (
	"context"
	"sync"

	"github.com/kakkky/hakoniwa/domain"
)

type agentRunner struct {
	state    agentState
	factory  agentFactory
	signalCh reconcileSignalCh
}

func newAgentRunner() *agentRunner {
	return &agentRunner{
		signalCh: make(reconcileSignalCh, 16),
		state: agentState{
			residentAgents: residentAgentsState{
				running: make(map[domain.ResidentID]context.CancelFunc),
			},
		},
	}
}

func (ar *agentRunner) setAgentFactory(factory agentFactory) {
	ar.factory = factory
}

type agentFactory struct {
	newGameMasterAgent func() *gameMasterAgent
	newResidentAgent   func(resident *domain.Resident) *residentAgent
}

type agentState struct {
	gameMasterAgent gameMasterAgentState
	residentAgents  residentAgentsState
}

type gameMasterAgentState struct {
	mu         sync.Mutex
	runnning   bool
	cancelFunc context.CancelFunc
}

type residentAgentsState struct {
	mu      sync.Mutex
	desired []residentAgent
	running map[domain.ResidentID]context.CancelFunc
}

type (
	reconcileSignal   struct{}
	reconcileSignalCh chan reconcileSignal
)

func (ar *agentRunner) run(ctx context.Context) error {
	gameMasterAgent := ar.factory.newGameMasterAgent()
	go gameMasterAgent.run(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ar.signalCh:
			ar.reconcileResidentAgents(ctx)
		}
	}
}

func (r *agentRunner) reconcileResidentAgents(ctx context.Context) {
	r.state.residentAgents.mu.Lock()
	defer r.state.residentAgents.mu.Unlock()
	for _, desired := range r.state.residentAgents.desired {
		if _, alreadyRunning := r.state.residentAgents.running[desired.residentID]; alreadyRunning {
			continue
		}
		agentCtx, cancelFn := context.WithCancel(ctx)

		residentAgent := &desired
		r.state.residentAgents.running[residentAgent.residentID] = cancelFn

		go residentAgent.run(agentCtx)
	}
}

func addResidentAgent(ar *agentRunner, resident *domain.Resident) *residentAgent {
	residentAgent := ar.factory.newResidentAgent(resident)
	ar.state.residentAgents.mu.Lock()
	ar.state.residentAgents.desired = append(ar.state.residentAgents.desired, *residentAgent)
	ar.state.residentAgents.mu.Unlock()

	ar.signalCh <- reconcileSignal{}
	return residentAgent
}
