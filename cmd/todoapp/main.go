package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/nurassul/todoapp-golang/internal/core/logger"
	"github.com/nurassul/todoapp-golang/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/nurassul/todoapp-golang/internal/core/transport/http/middleware"
	core_http_server "github.com/nurassul/todoapp-golang/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/nurassul/todoapp-golang/internal/features/statistics/repository/postgres"
	statistics_service "github.com/nurassul/todoapp-golang/internal/features/statistics/service"
	statistics_trasport_http "github.com/nurassul/todoapp-golang/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/nurassul/todoapp-golang/internal/features/tasks/repository/postgres"
	tasks_service "github.com/nurassul/todoapp-golang/internal/features/tasks/service"
	tasks_transport "github.com/nurassul/todoapp-golang/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/nurassul/todoapp-golang/internal/features/users/repository/postgres"
	users_service "github.com/nurassul/todoapp-golang/internal/features/users/service"
	users_transport_http "github.com/nurassul/todoapp-golang/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", timeZone))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("featrue", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_trasport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
