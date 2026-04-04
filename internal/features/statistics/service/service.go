package statistics_service

import (
	"context"
	"time"

	"github.com/Kor1992/todo/internal/core/domain"
)

type StatisticsService struct {
	repo StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(ctx context.Context, userID *int, from, to *time.Time) ([]domain.Task, error)
}

func NewStatisticsService(repo StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		repo: repo,
	}
}
