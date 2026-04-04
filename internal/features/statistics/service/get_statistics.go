package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Kor1992/todo/internal/core/domain"
	core_errors "github.com/Kor1992/todo/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(ctx context.Context, userID *int, from, to *time.Time) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("`to` must be after `from`: %w", core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.repo.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil

}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{
			TasksCreated:              0,
			TaskCompleted:             0,
			TasksCompletedRate:        nil,
			TasksAverageCompletedTime: nil,
		}
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var totalCompletionDuration time.Duration

	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++

		}
		compleationDiration := task.ComplitionDuration()
		if compleationDiration != nil {
			totalCompletionDuration += *compleationDiration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletedTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)
		tasksAverageCompletedTime = &avg
	}

	return domain.Statistics{
		TasksCreated:              tasksCreated,
		TaskCompleted:             tasksCompleted,
		TasksCompletedRate:        &tasksCompletedRate,
		TasksAverageCompletedTime: tasksAverageCompletedTime,
	}
}
