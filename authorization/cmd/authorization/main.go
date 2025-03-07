package main

import (
	"context"
	"github.com/megorka/goproject/authorization/internal/config"
	"github.com/megorka/goproject/authorization/internal/repository"
	"github.com/megorka/goproject/authorization/internal/service"
	router "github.com/megorka/goproject/authorization/internal/transport/http"
	"github.com/megorka/goproject/authorization/pkg/logger"
	"github.com/megorka/goproject/authorization/pkg/postgres"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	ctx, _ = logger.New(ctx)

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to read config", zap.Error(err))
	}
	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	repo := repository.NewRepository(db)

	svc := service.NewService(repo)

	handler := router.NewHandler(svc)

	r := router.NewRouter(cfg.Router, handler)

	go func() {
		err := r.Run(ctx)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to run server", zap.Error(err))
		}
	}()
	select {
	case <-ctx.Done():
		logger.GetLoggerFromCtx(ctx).Info(ctx, "shutting down")
		if err := r.Shutdown(ctx); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to shutdown server", zap.Error(err))
		}
		db.Close()
	}
}
