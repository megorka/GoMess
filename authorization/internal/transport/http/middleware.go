package router

import (
	"context"
	"github.com/google/uuid"
	"github.com/megorka/goproject/authorization/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

func MiddleWare(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, "request_id", requestID)

		var err error
		ctx, err = logger.New(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger := logger.GetLoggerFromCtx(ctx)
		logger.Info(ctx, "request started", zap.String("method", r.Method), zap.String("path", r.URL.Path))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
