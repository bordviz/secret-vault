package user

import (
	"context"
	"log/slog"
	"net/http"
	"vault/internal/config"

	"vault/pkg/handlers"
	"vault/pkg/lib/logger/sl"
	mwAuth "vault/pkg/lib/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	ErrNotFound = "vault not found"
)

type UserHandlerClient struct {
	userDBClient UserDB
	log          *slog.Logger
	secret       string
}

func AddUserRouter(r chi.Router, userClient UserDB, log *slog.Logger, cfg *config.Config) func(r chi.Router) {
	client := UserHandlerClient{
		userDBClient: userClient,
		log:          log,
		secret:       cfg.Secret,
	}

	return func(r chi.Router) {
		r.Use(mwAuth.UserAuth(log, cfg.Secret))

		r.Get("/get", client.GetVault(context.TODO()))
	}
}

func (h *UserHandlerClient) GetVault(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "user.handlers.GetVault"

		h.log = h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var key mwAuth.ContextKey = "vaultID"
		h.log.Debug("vault id on context value", slog.Any("vault_id", r.Context().Value(key)))

		id, ok := r.Context().Value(key).(int)
		if !ok {
			h.log.Error("failed to convert id to int")
			handlers.ErrorResponse(w, r, 500, "internal server error")
			return
		}

		vault, err := h.userDBClient.GetVault(ctx, h.log, id)
		if err != nil {
			if err.Error() == ErrNotFound {
				h.log.Error("vault not found", sl.OpErr(op, err))
				handlers.ErrorResponse(w, r, 404, err.Error())
				return
			}
			h.log.Error("failed to get vault", sl.OpErr(op, err))
			handlers.ErrorResponse(w, r, 500, err.Error())
			return
		}

		handlers.SuccessResponse(w, r, 200, vault)
		h.log.Info("vault successfully getted")
	}
}
