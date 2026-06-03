package domain

import "time"

type EventType string

const (
	EventTaskCreated EventType = "task.created"
	EventTaskDone    EventType = "task.done"
	EventTaskFailed  EventType = "task.failed"
)

type DomainEvent struct {
	Type      EventType
	TaskID    TaskID
	OccuredAt time.Time
}

func NewEvent(eventType EventType, taskID TaskID) DomainEvent {
	return DomainEvent{
		Type:      eventType,
		TaskID:    taskID,
		OccuredAt: time.Now(),
	}
}
