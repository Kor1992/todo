package users_transport_http

type UsersHTTPHandler struct {
	usersService UsersService
}
type UsersService interface {
}

func NewUsersHTTPHandler(userService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: userService,
	}
}
