package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_middleware "github.com/Kor1992/todo/internal/core/transport/http/middleware"
	core_http_server "github.com/Kor1992/todo/internal/core/transport/http/server"
	users_transport_http "github.com/Kor1992/todo/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(nil)
	log, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer log.Close()
	log.Debug("Starting ToDo app")

	usersRouters := usersTransportHTTP.Routers()

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersRouters...)

	httpServer := core_http_server.NewHTTPServer(core_http_server.NewConfigMust(),
		log,
		core_middleware.RequestId(),
		core_middleware.Logger(log),
		core_middleware.Panic(),
		core_middleware.Trace(),
	)

	httpServer.RegisterAPIRouter(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}

}
