package agents

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type Runtime struct {
	eventBroker       *eventBroker
	commandSubscriber *commandSubscriber
	agentRunner       *agentRunner
}

type AgentToolKit struct {
	ResidentRepository domain.ResidentRepository
}

func NewRuntime(llmProvider domain.LLMProvider, agentCommandCh domain.AgentCommandCh, toolKit AgentToolKit) *Runtime {
	// EventBroker
	eventBroker := newEventBroker()

	// AgentRunner
	agentRunner := newAgentRunner()
	sendEventFunc := func(e domain.Event) { eventBroker.inbox <- e }
	agentRunner.setAgentFactory(agentFactory{
		newResidentAgent: func(resident *domain.Resident) *residentAgent {
			return &residentAgent{
				agentBase:  newAgentBase(sendEventFunc, llmProvider),
				residentID: resident.ID,
				repository: toolKit.ResidentRepository,
			}
		},
		newGameMasterAgent: func() *gameMasterAgent {
			return &gameMasterAgent{
				agentBase: newAgentBase(sendEventFunc, llmProvider),
			}
		},
	})

	// CommandSubscriber
	commandSubscriber := newCommandSubscriber(agentCommandCh)
	registerHandler(
		commandSubscriber,
		domain.AddResidentAgentCommand{},
		func(ctx context.Context, cmd domain.AddResidentAgentCommand) error {
			newResidentAgent := addResidentAgent(agentRunner, &cmd.Resident)
			eventBroker.registerRoutes(cmd.Resident.ID, newResidentAgent.inbox)
			return nil
		})
	registerHandler(
		commandSubscriber,
		domain.PublishEventCommand{},
		func(ctx context.Context, cmd domain.PublishEventCommand) error {
			eventBroker.inbox <- cmd.Event
			return nil
		})

	return &Runtime{
		eventBroker:       eventBroker,
		commandSubscriber: commandSubscriber,
		agentRunner:       agentRunner,
	}
}

func (r *Runtime) Run(ctx context.Context) error {
	go r.eventBroker.run(ctx)
	go r.commandSubscriber.run(ctx)
	go r.agentRunner.run(ctx)

	return nil
}
