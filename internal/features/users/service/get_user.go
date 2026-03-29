package users_service

import (
	"context"
	"fmt"

	"github.com/Kor1992/todo/internal/core/domain"
	core_errors "github.com/Kor1992/todo/internal/core/errors"
)

func (s *UsersService) GetUser(ctx context.Context, id int) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %v : %w", err, core_errors.ErrNotFound)
	}

	return user, nil
}
