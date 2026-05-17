package agents

import (
	"context"
	"sync"

	"github.com/kakkky/hakoniwa/domain"
)

type agentEventInbox chan domain.Event

type brokerEventInbox = agentEventInbox

type eventBroker struct {
	inbox    brokerEventInbox
	routesMu sync.RWMutex
	routes   map[domain.ResidentID]agentEventInbox
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
	e.routesMu.Lock()
	defer e.routesMu.Unlock()
	e.routes[id] = inbox
}

func (e *eventBroker) dispatch(event domain.Event) {
	e.routesMu.RLock()
	to, ok := e.routes[event.To().ID]
	e.routesMu.RUnlock()
	if !ok {
		return
	}
	to <- event
}
