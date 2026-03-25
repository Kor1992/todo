package users_transport_http

import (
	"context"
	"net/http"

	"github.com/Kor1992/todo/internal/core/domain"
	core_http_server "github.com/Kor1992/todo/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}
type UsersService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
}

func NewUsersHTTPHandler(userService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: userService,
	}
}

func (h *UsersHTTPHandler) Routers() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
