package task_servis

import (
	"context"

	"github.com/Kor1992/todo/internal/core/domain"
)

type TasksService struct {
	repository TaskRepository
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userID, limit, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, taskID int) (domain.Task, error)
	DeleteTask(ctx context.Context, taskID int) error
	PatchTask(ctx context.Context, taskID int, task domain.Task) (domain.Task, error)
}

func NewTasksService(repo TaskRepository) *TasksService {
	return &TasksService{
		repository: repo,
	}
}
