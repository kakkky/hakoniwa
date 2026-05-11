package agents

type eventBroker struct {
	inbox  eventInbox
	routes map[ID]eventInbox
}

func newEventBroker() *eventBroker {
	return &eventBroker{}
}

func (e *eventBroker) run() error {
	return nil
}

func (e *eventBroker) registerRoutes(id ID, inbox eventInbox) error {
	return nil
}

func (e *eventBroker) dispatchMessage(msg message) error {
	return nil
}

func (e *eventBroker) dispatchOppotunity(op oppotunity) error {
	return nil
}
