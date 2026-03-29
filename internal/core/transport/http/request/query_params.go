package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/Kor1992/todo/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf("param=%s by key=%s not a valid int %v:%w", param, key, err, core_errors.ErrInvalidArgument)
	}
	return &val, nil
}

// func GetTimeQueryParam(r *http.Request, key string) int {

// }
