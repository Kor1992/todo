package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_postgres_pool "github.com/Kor1992/todo/internal/core/repositore/postgres/pool"
	core_middleware "github.com/Kor1992/todo/internal/core/transport/http/middleware"
	core_http_server "github.com/Kor1992/todo/internal/core/transport/http/server"
	users_postgres_repository "github.com/Kor1992/todo/internal/features/users/repository/postgres"
	users_service "github.com/Kor1992/todo/internal/features/users/service"
	users_transport_http "github.com/Kor1992/todo/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	log, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Debug("initializing pgx connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())

	if err != nil {
		log.Fatal("failes to init pgx connection pool", zap.Error(err))
	}
	defer pool.Close()
	log.Debug("initializing feature", zap.String("feature", "users"))

	// log.Debug("Starting ToDo app")
	usersRepository := users_postgres_repository.NewUserRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	log.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(core_http_server.NewConfigMust(),
		log,
		core_middleware.RequestId(),
		core_middleware.Logger(log),
		core_middleware.Panic(),
		core_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routers()...)
	httpServer.RegisterAPIRouter(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
