package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Kor1992/todo/internal/core/domain"
	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_http_request "github.com/Kor1992/todo/internal/core/transport/http/request"
	core_http_response "github.com/Kor1992/todo/internal/core/transport/http/response"
	core_http_types "github.com/Kor1992/todo/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Маким Максимович"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79998887766"`
}

func (p *PatchUserRequest) Validate() error {
	if p.FullName.Set {
		if p.FullName.Value == nil {
			return fmt.Errorf("fullName can't be null")
		}
		fullNameLen := len([]rune(*p.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("full name must be between 3 and 100 symbols")
		}
	}

	if p.PhoneNumber.Set {
		if p.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*p.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("phone number must be between 10 an 15 symbols")
			}

			if !strings.HasPrefix(*p.PhoneNumber.Value, "+") {
				return fmt.Errorf("phone number start with +")
			}
		}

	}
	return nil

}

type PatchUserResponse UserDtoResponse

// PatchUser     godoc
// @Summary      Изменение пользователя
// @Description  Изменение информации об уже существующем в системе пользователе
// @Description  ### Логика обновления полей (Three-state logic):
// @Description  1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description  2. **Явно передано значение**: `"phone_number": "+711122233344"` - устанавливает новый номер телефона в БД
// @Description  3. **Передан null**: `"phone_number": null` - очищает поле в БД (set to NULL)
// @Description  Ограничения: `full_name` не может быть выставлен как null
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id      path int          true            "ID изменяемого пользователя" Format(uuid)
// @Param        request body PatchUserRequest true            "PatchUser тело запроса"
// @Success      200 {object} PatchUserResponse                "Успешно изменённый пользователь"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      404 {object} core_http_response.ErrorResponse "User not found"
// @Failure      409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTEPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID patch value")
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndVolidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JsonResponse(response, http.StatusOK)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(request.FullName.ToDomain(), request.PhoneNumber.ToDomain())
}
