package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Kor1992/todo/docs"
	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_middleware "github.com/Kor1992/todo/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	log        *core_logger.Logger
	middleware []core_middleware.MiddleWare
}

func NewHTTPServer(config Config, log *core_logger.Logger, middleware ...core_middleware.MiddleWare) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouter(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router.WithMiddleware()))
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"GET /swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)

	s.mux.HandleFunc(
		"GET /swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_middleware.ChainMiddleWare(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()
	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shurdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shurdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server : %w", err)
		}
		s.log.Warn("HTTP server stoped")
	}
	return nil
}

func (s *HTTPServer) RegisterStatic() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Если запрос не на API и не на Swagger — отдаём index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./public/index.html")
			return
		}
		http.ServeFile(w, r, "./public/"+r.URL.Path)
	})
}
