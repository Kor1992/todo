package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_http_request "github.com/Kor1992/todo/internal/core/transport/http/request"
	core_http_response "github.com/Kor1992/todo/internal/core/transport/http/response"
)

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
