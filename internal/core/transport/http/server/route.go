package core_http_server

import (
	"net/http"

	core_middleware "github.com/Kor1992/todo/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_middleware.MiddleWare
}

func (r *Route) WithMiddleware() http.Handler {
	return core_middleware.ChainMiddleWare(r.Handler, r.Middleware...)
}
