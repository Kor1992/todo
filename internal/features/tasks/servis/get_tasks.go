package task_servis

import (
	"context"
	"fmt"

	"github.com/Kor1992/todo/internal/core/domain"
	core_errors "github.com/Kor1992/todo/internal/core/errors"
)

func (s *TasksService) GetTasks(ctx context.Context, userId, limit, offset *int) ([]domain.Task, error) {

	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit nust be non-negative %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset nust be non-negative %w", core_errors.ErrInvalidArgument)
	}
	users, err := s.repository.GetTasks(ctx, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository: %w", err)
	}
	return users, nil

}
