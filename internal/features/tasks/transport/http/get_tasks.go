package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_http_request "github.com/Kor1992/todo/internal/core/transport/http/request"
	core_http_response "github.com/Kor1992/todo/internal/core/transport/http/response"
)

type GetTasksReaponse []TasksDTOResponse

// GetTasks      godoc
// @Summary      Список задач
// @Description  Просмотр списка задач с опциональной пагинацией и/или фильтрацией по ID автора задачи
// @Tags         tasks
// @Produce      json
// @Param        user_id  query     int false  "Фильтрация задач по ID автора" Format(uuid)
// @Param        limit    query     int  false  "Размер страницы с задачами"
// @Param        offset   query     int  false  "Смещение страницы с задачами"
// @Success      200  {object}  []TasksDTOResponse  "Список задач"
// @Failure      400      {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      500      {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks [get]
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
