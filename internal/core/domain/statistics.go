package domain

import "time"

type Statistics struct {
	TasksCreated              int
	TaskCompleted             int
	TasksCompletedRate        *float64
	TasksAverageCompletedTime *time.Duration
}

func NewStatistics(TasksCreated int,
	TaskCompleted int,
	TasksCompletedRate *float64,
	TasksAverageCompletedTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:              TasksCreated,
		TaskCompleted:             TaskCompleted,
		TasksCompletedRate:        TasksCompletedRate,
		TasksAverageCompletedTime: TasksAverageCompletedTime,
	}
}
