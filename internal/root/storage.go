package root

import (
	"context"
	"log/slog"
	"vault/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=RootDB
type RootDB interface {
	CreateVault(ctx context.Context, log *slog.Logger, model models.SecretCreateDTO) (int, error)
	GetVault(ctx context.Context, log *slog.Logger, id int) (models.SecretModel, error)
	CheckVault(ctx context.Context, log *slog.Logger, id int) error
}
