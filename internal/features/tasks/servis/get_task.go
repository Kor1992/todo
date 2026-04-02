package task_servis

import (
	"context"
	"fmt"

	"github.com/Kor1992/todo/internal/core/domain"
)

func (s *TasksService) GetTask(ctx context.Context, taskID int) (domain.Task, error) {
	task, err := s.repository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	return task, nil
}
