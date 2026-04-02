package task_servis

import (
	"context"
	"fmt"

	"github.com/Kor1992/todo/internal/core/domain"
)

func (s *TasksService) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	domainTask, err := s.repository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return domainTask, nil

}
