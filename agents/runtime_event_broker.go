package agents

import (
	"context"

	"github.com/kakkky/hakoniwa/domain"
)

type agentEventInbox chan domain.Event

type brokerEventInbox = agentEventInbox

type eventBroker struct {
	inbox  brokerEventInbox
	routes map[domain.ResidentID]agentEventInbox
}

func newEventBroker() *eventBroker {
	inbox := make(chan domain.Event, 32)
	routes := make(map[domain.ResidentID]agentEventInbox)
	return &eventBroker{
		inbox:  inbox,
		routes: routes,
	}
}

func (e *eventBroker) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-e.inbox:
			e.dispatch(event)
		}
	}
}

func (e *eventBroker) registerRoutes(id domain.ResidentID, inbox agentEventInbox) {
	e.routes[id] = inbox
}

func (e *eventBroker) dispatch(event domain.Event) {
	to := e.routes[event.To().ID]
	to <- event
}
