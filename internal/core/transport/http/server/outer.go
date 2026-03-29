package core_http_server

import (
	"fmt"
	"net/http"

	core_middleware "github.com/Kor1992/todo/internal/core/transport/http/middleware"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	middleware []core_middleware.MiddleWare
}

func NewApiVersionRouter(apiVersion ApiVersion, middleware ...core_middleware.MiddleWare) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, router := range routes {
		pattern := fmt.Sprintf("%s %s", router.Method, router.Path)

		r.Handle(pattern, router.WithMiddleware())
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_middleware.ChainMiddleWare(r, r.middleware...)
}
