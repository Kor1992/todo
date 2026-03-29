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

type PatcUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (p *PatcUserRequest) Validate() error {
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

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTEPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID patch value")
		return
	}

	var request PatcUserRequest
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

func userPatchFromRequest(request PatcUserRequest) domain.UserPatch {
	return domain.NewUserPatch(request.FullName.ToDomain(), request.PhoneNumber.ToDomain())
}
