package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TaskID string

func NewTaskID() TaskID {
	return TaskID(uuid.New().String())
}

type Status string

const (
	StatusPending  Status = "pending"
	StatusAssigned Status = "assigned"
	StatusRunning  Status = "running"
	StatusDone     Status = "done"
	StatusFailed   Status = "failed"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrTaskFailed = errors.New("task failed")

func (s Status) String() string { return string(s) }

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func ParsePriority(s string) (Priority, error) {
	switch strings.ToLower(s) {
	case "low":
		return PriorityLow, nil
	case "medium":
		return PriorityMedium, nil
	case "high":
		return PriorityHigh, nil
	default:
		return "", fmt.Errorf("invalid priority %q: must be low, medium, or high", s)
	}
}

// --- Entity ---

type Task struct {
	ID        TaskID
	Title     string
	Priority  Priority
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateTaskDTO struct {
	Title    *string
	Priority *string
	Status   *Status
}

func NewTask(title string, priority Priority) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("task title cannot be empty")
	}
	now := time.Now()
	return &Task{
		ID:        NewTaskID(),
		Title:     title,
		Priority:  priority,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (t *Task) MarkRunning() {
	t.Status = StatusRunning
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkDone() {
	t.Status = StatusDone
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkFailed() {
	t.Status = StatusFailed
	t.UpdatedAt = time.Now()
}

func (t *Task) String() string {
	return fmt.Sprintf("[%s] %s (priority: %s, status: %s)",
		t.ID, t.Title, t.Priority, t.Status)
}
