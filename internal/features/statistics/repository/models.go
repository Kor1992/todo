package statistics_postgres_repository

import (
	"time"

	"github.com/Kor1992/todo/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func tasksModelsFromDomain(models []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(models))

	for i, model := range models {
		domains[i] = taskDomainFromModel(model)
	}

	return domains
}

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}
