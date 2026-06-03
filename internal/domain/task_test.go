package domain_test

import (
	"taskflow/internal/domain"
	"testing"
)

func TestNewTask_Valid(t *testing.T) {
	task, err := domain.NewTask("send email", domain.PriorityHigh)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.Status != domain.StatusPending {
		t.Errorf("expected pending, got %s", task.Status)
	}
	if task.ID == "" {
		t.Error("expected a non-empty ID")
	}
}

func TestNewTask_EmptyTitle(t *testing.T) {
	_, err := domain.NewTask("", domain.PriorityLow)
	if err == nil {
		t.Fatal("expected an error for empty title")
	}
}

func TestParsePriority_Invalid(t *testing.T) {
	_, err := domain.ParsePriority("urgent")
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
}

func TestTask_StatusTransitions(t *testing.T) {
	task, _ := domain.NewTask("process invoice", domain.PriorityMedium)

	task.MarkRunning()
	if task.Status != domain.StatusRunning {
		t.Errorf("expected running, got %s", task.Status)
	}

	task.MarkDone()
	if task.Status != domain.StatusDone {
		t.Errorf("expected done, got %s", task.Status)
	}
}
