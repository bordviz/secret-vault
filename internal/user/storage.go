package user

import (
	"context"
	"log/slog"
	"vault/internal/models"
)

type UserDB interface {
	GetVault(ctx context.Context, log *slog.Logger, id int) (models.SecretModel, error)
}
