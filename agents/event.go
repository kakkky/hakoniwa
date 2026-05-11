package agents

type anyEvent interface {
}

type event[T message | oppotunity] struct {
	eventType eventType
	to        ID
	from      ID
	payroad   T
}

type eventType int

const (
	Unspesified     eventType = iota
	MessageEvent              = iota
	OppotunityEvent           = iota
)

type message struct{}

type oppotunity struct{}

type eventInbox chan<- anyEvent
