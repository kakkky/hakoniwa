package agents

type agentEvent interface {
	to() (id, name)
	from() (id, name)
	payload() string
}

type eventBase struct {
	toID     id
	toName   name
	fromID   id
	fromName name
}

func (eb eventBase) to() (id id, name name) {
	return eb.toID, eb.toName
}

func (eb eventBase) from() (id id, name name) {
	return eb.fromID, eb.fromName
}

type messageEvent struct {
	eventBase
	message string
}

func (me messageEvent) payload() string {
	return me.message
}

type opportunityEvent struct {
	eventBase
	opportunity string
}

func (oe opportunityEvent) payload() string {
	return oe.opportunity
}
