package users_postgres_repository

import core_postgres_pool "github.com/Kor1992/todo/internal/core/repositore/postgres/pool"

type UserRepository struct {
	pool core_postgres_pool.Pool
}

func NewUserRepository(pool core_postgres_pool.Pool) *UserRepository {

	return &UserRepository{
		pool: pool,
	}

}
