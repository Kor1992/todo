package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_http_request "github.com/Kor1992/todo/internal/core/transport/http/request"
	core_http_response "github.com/Kor1992/todo/internal/core/transport/http/response"
)

type GetTasksReaponse []TasksDTOResponse

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTEPResponseHandler(log, rw)

	userID, limit, offset, err := GetUserIdLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID/limit/offset query param")
		return
	}

	tasksDomain, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := GetTasksReaponse(taskDTOsFromDomains(tasksDomain))

	responseHandler.JsonResponse(response, http.StatusOK)

}

func GetUserIdLimitOffsetQueryParam(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDquetyParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDquetyParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user id query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get limit query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get offset query param: %w", err)
	}
	return userID, limit, offset, nil
}
