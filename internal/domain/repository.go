package domain

import "context"

type TaskRepository interface {
	Save(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id TaskID) (*Task, error)
	FindAll(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
}
