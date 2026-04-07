package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/Kor1992/todo/internal/core/config"
	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_pgx_pool "github.com/Kor1992/todo/internal/core/repositore/postgres/pool/pgx"
	core_middleware "github.com/Kor1992/todo/internal/core/transport/http/middleware"
	core_http_server "github.com/Kor1992/todo/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/Kor1992/todo/internal/features/statistics/repository"
	statistics_service "github.com/Kor1992/todo/internal/features/statistics/service"
	statistics_transport_http "github.com/Kor1992/todo/internal/features/statistics/transport/http"
	task_postgres "github.com/Kor1992/todo/internal/features/tasks/repository/postgres"
	task_servis "github.com/Kor1992/todo/internal/features/tasks/servis"
	tasks_transport_http "github.com/Kor1992/todo/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Kor1992/todo/internal/features/users/repository/postgres"
	users_service "github.com/Kor1992/todo/internal/features/users/service"
	users_transport_http "github.com/Kor1992/todo/internal/features/users/transport/http"
	"go.uber.org/zap"

	_ "github.com/Kor1992/todo/docs"
)

// @title        Golang Todo API
// @version      1.0
// @description  Todo Application REST-API scheme
// @host         127.0.0.1:5050
// @BasePath     /api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

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

	log.Debug("application time zone", zap.Any("zone", time.Local))

	log.Debug("initializing pgx connection pool")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())

	if err != nil {
		log.Fatal("failes to init pgx connection pool", zap.Error(err))
	}
	defer pool.Close()
	log.Debug("initializing feature", zap.String("feature", "users"))

	// log.Debug("Starting ToDo app")
	usersRepository := users_postgres_repository.NewUserRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	log.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepositore := task_postgres.NewTasksRepository(pool)
	tasksService := task_servis.NewTasksService(tasksRepositore)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	log.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	log.Debug("initializing HTTP server")
	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		log,
		core_middleware.CORS(),
		core_middleware.RequestId(),
		core_middleware.Logger(log),
		core_middleware.Trace(),
		core_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routers()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	// apiVersionRouterV2 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion2, core_middleware.Dummy("api v2 middleware"))
	// apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routers()...)

	httpServer.RegisterAPIRouter(apiVersionRouterV1) // apiVersionRouterV2,

	httpServer.RegisterSwagger()
	httpServer.RegisterStatic()

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
