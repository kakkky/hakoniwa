package domain

// ActorID はイベントの送受信主体 (住人 / 管理人 等) を識別する ID
type ActorID string

// ActorName はイベントの送受信主体 (住人 / 管理人 等) の表示名
type ActorName string

const (
	BuildManagerActorID   ActorID   = "manager"
	BuoldManagerActorName ActorName = "管理人"
	GameMasterActorID     ActorID   = "game_master"
	GameMasterActorName   ActorName = "ゲームマスター"
)

type Event interface {
	To() EventTo
	From() EventFrom
	Payload() string
}

type EventBase struct {
	EventTo   EventTo
	EventFrom EventFrom
}

type EventTo struct {
	ID   ResidentID
	Name ResidentName
}

type EventFrom struct {
	ID   ActorID
	Name ActorName
}

func (eb EventBase) To() EventTo     { return eb.EventTo }
func (eb EventBase) From() EventFrom { return eb.EventFrom }

type MessageEvent struct {
	EventBase
	Message string
}

func (me MessageEvent) Payload() string { return me.Message }

type OpportunityEvent struct {
	EventBase
	Opportunity string
}

func (oe OpportunityEvent) Payload() string { return oe.Opportunity }
