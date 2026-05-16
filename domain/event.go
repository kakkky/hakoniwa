package domain

import "time"

type Event interface {
	_isEvent()
}

type EventInbox chan Event

// WorldEvent: 世界全体の出来事 (Tick由来など)
type WorldEvent struct {
	Content    string
	OccurredAt time.Time
}

func (WorldEvent) _isEvent() {}

func NewWorldEvent(content string, now time.Time) WorldEvent {
	return WorldEvent{Content: content, OccurredAt: now}
}

func NewTickMorningEvent(now time.Time) WorldEvent {
	return NewWorldEvent("世界は朝になった", now)
}

func NewTickNightEvent(now time.Time) WorldEvent {
	return NewWorldEvent("世界は夜になった", now)
}

// ResidentEvent: 個別residentの出来事
type ResidentEvent struct {
	ResidentID ResidentID
	Content    string
	OccurredAt time.Time
}

func (ResidentEvent) _isEvent() {}

func NewResidentEvent(id ResidentID, content string, now time.Time) ResidentEvent {
	return ResidentEvent{ResidentID: id, Content: content, OccurredAt: now}
}
