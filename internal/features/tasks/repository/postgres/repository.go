package task_postgres

import core_postgres_pool "github.com/Kor1992/todo/internal/core/repositore/postgres/pool"

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(pool core_postgres_pool.Pool) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
