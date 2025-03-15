package router

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/megorka/goproject/post_service/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "Authorization header is required")
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := tokenParts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "Unexpected signing method: ")
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "Failed to parse token")
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "Invalid user_id in token")
			http.Error(w, "Invalid user_id in token", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user_id", userID))
		next.ServeHTTP(w, r)
	})
}

func Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, "request_id", requestID)

		var err error
		ctx, err = logger.New(ctx)
		if err != nil {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "Failed to create logger", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger := logger.GetLoggerFromCtx(ctx)
		logger.Info(ctx, "request started", zap.String("method", r.Method), zap.String("path", r.URL.Path))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
