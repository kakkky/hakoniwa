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

func NewRuntime(llmProvider domain.LLMProvider) *Runtime {
	// EventBroker
	eventBroker := newEventBroker()

	// AgentRunner
	agentRunner := newAgentRunner()
	agentBase := newAgentBase(
		func(e domain.Event) { eventBroker.inbox <- e },
		llmProvider,
	)
	agentRunner.setAgentFactory(agentFactory{
		newResidentAgent: func(resident *domain.Resident) *residentAgent {
			return &residentAgent{
				agentBase: agentBase,
				resident:  resident,
			}
		},
		newGameMasterAgent: func() *gameMasterAgent {
			return &gameMasterAgent{
				agentBase: agentBase,
			}
		},
	})

	// CommandSubscriber
	commandSubscriber := newCommandSubscriber()
	registerHandler(
		commandSubscriber,
		domain.AddResidentAgentCommand{},
		func(ctx context.Context, cmd domain.AddResidentAgentCommand) error {
			addResidentAgent(agentRunner, &cmd.Resident)
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

// CommandInbox は外部から command を流し込むためのチャネルを返す。
// AgentCommander の実装にこのチャネルを渡して利用する。
func AgentCommandInbox(r *Runtime) domain.AgentCommandInbox {
	return r.commandSubscriber.cmdInbox
}
