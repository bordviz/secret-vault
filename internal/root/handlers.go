package root

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"vault/internal/config"
	"vault/internal/models"
	"vault/pkg/handlers"
	"vault/pkg/lib/jwt"
	"vault/pkg/lib/logger/sl"
	mwAuth "vault/pkg/lib/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var (
	ErrNotFound = "vault not found"
)

type RootHandlerClient struct {
	rootDBClient RootDB
	log          *slog.Logger
	secret       string
}

func AddRootRouter(r chi.Router, rootClient RootDB, log *slog.Logger, cfg *config.Config) func(r chi.Router) {
	client := NewRootHandlerClient(rootClient, log, cfg.Secret)

	return func(r chi.Router) {
		r.Use(mwAuth.RootAuth(log, cfg.RootToken))

		r.Post("/create", client.CreateVault(context.TODO()))
		r.Get("/get/{id}", client.GetVault(context.TODO()))
		r.Post("/create-token", client.CreateVaultToken(context.TODO()))
	}
}

func NewRootHandlerClient(rootClient RootDB, log *slog.Logger, secret string) *RootHandlerClient {
	return &RootHandlerClient{
		rootDBClient: rootClient,
		log:          log,
		secret:       secret,
	}
}

func (h *RootHandlerClient) CreateVault(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "root.handlers.CreateVault"

		h.log = h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var model models.SecretCreateModel
		err := render.DecodeJSON(r.Body, &model)
		if err != nil {
			h.log.Error("failed to decode model", sl.OpErr(op, err))
			handlers.ErrorResponse(w, r, 400, "failed to decode model")
			return
		}
		if err := model.Validate(); err != nil {
			h.log.Error("validate error", sl.OpErr(op, err))
			handlers.ErrorResponse(w, r, 422, err.Error())
			return
		}

		//TODO: add data encode

		id, err := h.rootDBClient.CreateVault(ctx, h.log, model.ConvertToDTO())
		if err != nil {
			h.log.Error("failed to save new vault on database", sl.OpErr(op, err))
			handlers.ErrorResponse(w, r, 500, err.Error())
			return
		}

		handlers.SuccessResponse(w, r, 201, map[string]any{
			"message": "new vault successfully created",
			"id":      id,
		})
		h.log.Info("new vault successfully created", "id", id)
	}
}

func (h *RootHandlerClient) GetVault(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "root.handlers.GetVault"

		h.log = h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			h.log.Error("query param id is empty")
			handlers.ErrorResponse(w, r, 400, "query param id is empty")
			return
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			h.log.Error("failed to convert id to integer", sl.OpErr(op, err))
			handlers.ErrorResponse(w, r, 400, "query parameter must be int")
			return
		}
		model, err := h.rootDBClient.GetVault(ctx, h.log, idInt)
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

		//TODO: add data decode

		handlers.SuccessResponse(w, r, 200, model)
		h.log.Info("secret vault successfully getted")
	}
}

func (h *RootHandlerClient) CreateVaultToken(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "root.handlers.CreateVaultToken"

		h.log = h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var model models.CreateVaultTokenDTO

		if err := render.Decode(r, &model); err != nil {
			h.log.Error("failed to decode model", sl.OpErr(op, err))
			handlers.ErrorResponse(w, r, 400, "failed to decode model")
			return
		}

		if err := model.Validate(); err != nil {
			h.log.Error("validate error", sl.OpErr(op, err))
			handlers.ErrorResponse(w, r, 422, err.Error())
			return
		}

		checkVault := h.rootDBClient.CheckVault(ctx, h.log, model.VaultID)
		if checkVault != nil {
			h.log.Error("failed to check vault", sl.Err(checkVault))
			if checkVault.Error() == ErrNotFound {
				h.log.Error("vault not found", sl.Err(checkVault))
				handlers.ErrorResponse(w, r, 404, checkVault.Error())
				return
			}
			h.log.Error("failed to check vault", sl.OpErr(op, checkVault))
			handlers.ErrorResponse(w, r, 500, checkVault.Error())
			return
		}

		token, err := jwt.CreateToken(model.VaultID, h.secret, model.Expires)
		if err != nil {
			h.log.Error("failed to create new token", sl.OpErr(op, err))
			handlers.ErrorResponse(w, r, 500, "failed to create new token")
			return
		}

		handlers.SuccessResponse(w, r, 201,
			map[string]string{"token": token},
		)
		h.log.Info("vault token successfully created")
	}
}
