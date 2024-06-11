package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"vault/pkg/handlers"
	"vault/pkg/lib/jwt"
)

type ContextKey string

func RootAuth(log *slog.Logger, rootToken string) func(next http.Handler) http.Handler {
	const op = "middleware.auth.RootAuth"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Debug("recovered", slog.String("op", op), slog.String("detail", fmt.Sprintf("%s", rec)))
					log.Error("unauthorized", slog.String("op", op))
					handlers.ErrorResponse(w, r, 401, "unauthorized")
				}
			}()

			token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
			if token != rootToken {
				log.Error("unauthorized", slog.String("op", op), slog.String("token", token))
				handlers.ErrorResponse(w, r, 401, "unauthorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserAuth(log *slog.Logger, secret string) func(next http.Handler) http.Handler {
	const op = "middleware.auth.UserAuth"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Debug("recovered", slog.String("op", op), slog.String("detail", fmt.Sprintf("%s", rec)))
					log.Error("unauthorized", slog.String("op", op))
					handlers.ErrorResponse(w, r, 401, "unauthorized")
				}
			}()

			token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
			if token == "" {
				log.Error("unauthorized", slog.String("op", op), slog.String("token", token))
				handlers.ErrorResponse(w, r, 401, "unauthorized")
				return
			}

			id, err := jwt.DecodeToken(token, secret)
			log.Debug("vault id from token", "id", id)

			if err != nil {
				log.Error("failed to decode token", slog.String("op", op), slog.String("token", token))
				handlers.ErrorResponse(w, r, 401, "unauthorized")
			}

			var key ContextKey = "vaultID"

			ctx := context.WithValue(r.Context(), key, id)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
