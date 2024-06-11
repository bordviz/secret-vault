package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"vault/pkg/lib/logger/sl"
	"vault/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	Attemps  int
	Timeout  time.Duration
	Delay    time.Duration
}

type Client interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, log *slog.Logger, cfg DatabaseConfig) (pool *pgxpool.Pool) {
	const op = "database.postgresql.NewClient"

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	log.Debug("database dsn", "dsn", dsn)

	err := utils.DoWithTries(func() error {
		log.Debug("database connection attempt")
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		pool, _ = pgxpool.New(ctx, dsn)
		err := pool.Ping(ctx)

		if err != nil {
			log.Debug("database conection failed", sl.OpErr(op, err))
		}

		return err
	}, cfg.Attemps, cfg.Delay)

	if err != nil {
		log.Error("database conection failed", sl.Err(err))
		return nil
	}

	return pool
}
