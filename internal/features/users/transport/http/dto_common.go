package users_transport_http

import "github.com/Kor1992/todo/internal/core/domain"

type UserDtoResponse struct {
	ID          int     `json:"id" example:"10"`
	Version     int     `json:"version" example:"3"`
	FullName    string  `json:"full_name"  example:"Иван Петрович"`
	PhoneNumber *string `json:"phone_number" example:"+79998887766"`
}

func userDTOFromDomain(user domain.User) UserDtoResponse {
	return UserDtoResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDtoResponse {
	usersDTO := make([]UserDtoResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}
	return usersDTO
}
