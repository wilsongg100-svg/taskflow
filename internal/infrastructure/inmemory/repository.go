package inmemory

import (
	"context"
	"sync"
	"taskflow/internal/domain"
)

type Repository struct {
	mu    sync.RWMutex
	tasks map[domain.TaskID]*domain.Task
}

func NewRepository() *Repository {
	return &Repository{
		tasks: make(map[domain.TaskID]*domain.Task),
	}
}

func (r *Repository) Save(ctx context.Context, task *domain.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	task, exists := r.tasks[id]
	if !exists {
		return nil, domain.ErrTaskNotFound
	}
	taskCopy := *task
	return &taskCopy, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]*domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) Update(ctx context.Context, task *domain.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; !exists {
		return domain.ErrTaskNotFound
	}
	r.tasks[task.ID] = task
	return nil
}
