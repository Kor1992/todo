package core_middleware

import "net/http"

type MiddleWare func(http.Handler) http.Handler
