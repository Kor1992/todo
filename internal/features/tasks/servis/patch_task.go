package task_servis

import (
	"context"
	"fmt"

	"github.com/Kor1992/todo/internal/core/domain"
)

func (s *TasksService) PatchTask(ctx context.Context, taskID int, patch domain.TaskPatch) (domain.Task, error) {
	task, err := s.repository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task: %w", err)
	}

	domainTask, err := s.repository.PatchTask(ctx, taskID, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return domainTask, nil

}
