package statistics_transport_http

import "github.com/Kor1992/todo/internal/core/domain"

type GetStatisticsDto struct {
	TasksCreated              int      `json:"tasks_created"`
	TaskCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate        *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletedTime *string  `json:"tasks_averege_completed_time"`
}

func toDTOFromDomain(stat domain.Statistics) GetStatisticsDto {
	var avgTime *string
	if stat.TasksAverageCompletedTime != nil {
		duration := stat.TasksAverageCompletedTime.String()
		avgTime = &duration
	}

	return GetStatisticsDto{
		TasksCreated:              stat.TasksCreated,
		TaskCompleted:             stat.TaskCompleted,
		TasksCompletedRate:        stat.TasksCompletedRate,
		TasksAverageCompletedTime: avgTime,
	}
}
