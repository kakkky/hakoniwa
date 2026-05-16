package agents

import "context"

type agentEventInbox chan agentEvent

type brokerEventInbox = agentEventInbox

type eventBroker struct {
	inbox  brokerEventInbox
	routes map[id]agentEventInbox
}

func newEventBroker() *eventBroker {
	inbox := make(chan agentEvent, 32)
	routes := make(map[id]agentEventInbox)
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

func (e *eventBroker) registerRoutes(id id, inbox agentEventInbox) {
	e.routes[id] = inbox
}

func (e *eventBroker) dispatch(event agentEvent) {
	to := e.routes[event.to()]
	to <- event
}
