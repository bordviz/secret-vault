package root

import (
	"context"
	"log/slog"
	"vault/internal/models"
)

type RootDB interface {
	CreateVault(ctx context.Context, log *slog.Logger, model models.SecretCreateDTO) (int, error)
	GetVault(ctx context.Context, log *slog.Logger, id int) (models.SecretModel, error)
	CheckVault(ctx context.Context, log *slog.Logger, id int) error
}
