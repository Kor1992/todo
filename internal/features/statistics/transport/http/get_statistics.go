package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Kor1992/todo/internal/core/domain"
	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_http_request "github.com/Kor1992/todo/internal/core/transport/http/request"
	core_http_response "github.com/Kor1992/todo/internal/core/transport/http/response"
)

type GetStatisticsDto struct {
	TasksCreated              int      `json:"tasks_created"`
	TaskCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate        *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletedTime *string  `json:"tasks_averege_completed_time"`
}

func toDTOFromDomain(stat domain.Statistics) GetStatisticsDto {
	var avgTime *string
	if stat.TasksAverageCompletedTime != nil {
		duration := stat.TasksAverageCompletedTime.String()
		avgTime = &duration
	}

	return GetStatisticsDto{
		TasksCreated:              stat.TasksCreated,
		TaskCompleted:             stat.TaskCompleted,
		TasksCompletedRate:        stat.TasksCompletedRate,
		TasksAverageCompletedTime: avgTime,
	}
}

// GetStatistics godoc
// @Summary      Получение статистики
// @Description  Получение статистики по задачам с опциональной фильтрацией по user_id и/или временному промежутку
// @Tags         statistics
// @Produce      json
// @Param        user_id  query     int     false  "Фильтрация статистики по конкретному пользователю"
// @Param        from     query     string  false "Начало промежутка рассмотрения статистики (включительно), формат: YYYY-MM-DD"
// @Param        to       query     string  false "Конец промежутся рассмотрения статистики (не включительно), формат: YYYY-MM-DD"
// @Success      200      {object}  GetStatisticsDto "Успешное получение статистики"
// @Failure      400      {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      500      {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTEPResponseHandler(log, rw)

	userId, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "get query param")
		return
	}

	statisticsDomain, err := h.service.GetStatistics(ctx, userId, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	response := toDTOFromDomain(statisticsDomain)

	responseHandler.JsonResponse(response, http.StatusOK)

}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user id query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("get from query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("get to query param: %w", err)
	}

	return userID, from, to, nil

}
