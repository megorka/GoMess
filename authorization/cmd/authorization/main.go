package main

import (
	"context"
	"github.com/megorka/goproject/authorization/internal/repository"
	"github.com/megorka/goproject/authorization/internal/service"
	router "github.com/megorka/goproject/authorization/internal/transport/http"
	"github.com/megorka/goproject/authorization/pkg/logger"
	"go.uber.org/zap"

	"github.com/megorka/goproject/authorization/internal/config"
	"github.com/megorka/goproject/authorization/pkg/postgres"
)

func main() {

	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to read config", zap.Error(err))
	}

	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to database", zap.Error(err))
	}

	repo := repository.NewRepository(db)

	svc := service.NewService(repo)

	handler := router.NewHandler(svc)

	r := router.NewRouter(cfg.Router, handler)

	r.Run()
}
