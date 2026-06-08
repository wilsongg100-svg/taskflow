package application

import (
	"context"
	"fmt"
	"log"
	"taskflow/internal/domain"
)

type TaskService struct {
	repo   domain.TaskRepository
	events chan domain.DomainEvent
}

func NewTaskService(repo domain.TaskRepository, events chan domain.DomainEvent) *TaskService {
	return &TaskService{
		repo:   repo,
		events: events,
	}
}

func (t *TaskService) Create(ctx context.Context, title string, priority string) (*domain.Task, error) {
	parsedPriority, err := domain.ParsePriority(priority)
	if err != nil {
		return nil, err
	}

	newTask, err := domain.NewTask(title, parsedPriority)
	if err != nil {
		return nil, err
	}

	if err := t.repo.Save(ctx, newTask); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	if sendEventErr := t.sendEvent(ctx, domain.NewEvent(domain.EventTaskCreated, newTask.ID)); sendEventErr != nil {
		log.Printf("[ERROR] Failed to send event for task %s: %v", newTask.ID, sendEventErr)
	}

	return newTask, nil
}

func (t *TaskService) List(ctx context.Context) ([]*domain.Task, error) {
	return t.repo.FindAll(ctx)
}

func (t *TaskService) Get(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	return t.repo.FindByID(ctx, id)
}

func (t *TaskService) Update(ctx context.Context, id domain.TaskID, dto domain.UpdateTaskDTO) error {
	currentTask, err := t.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding task: %w", err)
	}

	updatedTask := *currentTask

	if dto.Title != nil {
		updatedTask.Title = *dto.Title
	}
	if dto.Priority != nil {
		parsedPriority, err := domain.ParsePriority(*dto.Priority)
		if err != nil {
			return fmt.Errorf("invalid task priority: %w", err)
		}
		updatedTask.Priority = parsedPriority
	}

	if dto.Status != nil {
		updatedTask.Status = *dto.Status
	}

	if err := t.repo.Update(ctx, &updatedTask); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	if sendEventErr := t.sendEvent(ctx, domain.NewEvent(domain.EventTaskUpdated, currentTask.ID)); sendEventErr != nil {
		log.Printf("[ERROR] Failed to send event for task %s: %v", currentTask.ID, sendEventErr)
	}
	return nil
}

func (t *TaskService) MarkDone(ctx context.Context, id domain.TaskID) error {
	task, err := t.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding task: %w", err)
	}
	task.MarkDone()

	if err := t.repo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to mark task as done: %w", err)
	}
	if sendEventErr := t.sendEvent(ctx, domain.NewEvent(domain.EventTaskDone, task.ID)); sendEventErr != nil {
		log.Printf("[ERROR] Failed to send event for task %s: %v", task.ID, sendEventErr)
	}
	return nil
}

func (t *TaskService) MarkFailed(ctx context.Context, id domain.TaskID) error {
	task, err := t.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding task: %w", err)
	}
	task.MarkFailed()

	if err := t.repo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to mark task as failed: %w", err)
	}
	if sendEventErr := t.sendEvent(ctx, domain.NewEvent(domain.EventTaskFailed, task.ID)); sendEventErr != nil {
		log.Printf("[ERROR] Failed to send event for task %s: %v", task.ID, sendEventErr)
	}
	return nil
}

func (t *TaskService) sendEvent(ctx context.Context, event domain.DomainEvent) error {
	select {
	case t.events <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		log.Printf("[WARN] Event buffer full. Dropped or deferred event for task: %s", event.TaskID)
		return nil
	}
}
