package db

import (
	"context"
	"errors"
	"log/slog"
	"vault/internal/models"
	"vault/pkg/database/postgresql"
	"vault/pkg/lib/logger/sl"
	"vault/pkg/utils"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBClient struct {
	dbClient postgresql.Client
}

func NewClient(dbClient postgresql.Client) *DBClient {
	return &DBClient{
		dbClient: dbClient,
	}
}

func (r *DBClient) CreateVault(ctx context.Context, log *slog.Logger, model models.SecretCreateDTO) (int, error) {
	const op = "db.postgresql.CreateVault"

	tx, err := r.dbClient.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", sl.OpErr(op, err))
		return 0, errors.New("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	createVaultQuery := `
		INSERT INTO vault
			(name)
		VALUES 
			($1)
		RETURNING id;
	`

	log.Debug("create vault query", slog.String("op", op), slog.String("query", utils.QueryConvert(createVaultQuery)))

	var id int

	if err := tx.QueryRow(context.TODO(), createVaultQuery, model.Name).Scan(&id); err != nil {
		log.Error("failed to create new vault", sl.OpErr(op, err))
		tx.Rollback(ctx)
		return 0, errors.New("failed to create new vault")
	}

	createValueQuery := `
		INSERT INTO value
			(vault_id, key, value)
		VALUES ($1, $2, $3);
	`

	log.Debug("create value query", slog.String("op", op), slog.String("query", utils.QueryConvert(createValueQuery)))

	rows := make([][]interface{}, len(model.Data))
	for i, v := range model.Data {
		rows[i] = []interface{}{id, v.Key, v.Value}
	}

	log.Debug("print rows", "rows", rows)

	for _, row := range rows {
		_, err := tx.Exec(ctx, createValueQuery, row...)
		if err != nil {
			log.Error("failed to insert values to database", sl.OpErr(op, err))
			tx.Rollback(ctx)
			return 0, errors.New("failed to insert values to database")
		}
	}
	err = tx.Commit(ctx)

	return id, err
}

func (r *DBClient) GetVault(ctx context.Context, log *slog.Logger, id int) (models.SecretModel, error) {
	const op = "db.postgresql.GetVault"

	tx, err := r.dbClient.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", sl.OpErr(op, err))
		return models.SecretModel{}, errors.New("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	getVaultQuery := `
		SELECT id, name FROM vault
		WHERE id = $1;
	`

	log.Debug("get vault query", slog.String("op", op), slog.String("query", utils.QueryConvert(getVaultQuery)))

	var vault models.VaultModel
	err = tx.QueryRow(ctx, getVaultQuery, id).Scan(&vault.ID, &vault.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error("failed to get vault", sl.OpErr(op, err))
			return models.SecretModel{}, errors.New("vault not found")
		}
		log.Error("failed to get vault", sl.OpErr(op, err))
		return models.SecretModel{}, errors.New("failed to get vault")
	}

	getValuesQuery := `
		SELECT key, value FROM value
		WHERE vault_id = $1;
	`

	log.Debug("get values query", slog.String("op", op), slog.String("query", utils.QueryConvert(getValuesQuery)))

	rows, err := tx.Query(ctx, getValuesQuery, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsNoData(pgErr.Error()) {
			log.Error("failed to get vault values", sl.OpErr(op, err))
			return models.SecretModel{}, errors.New("failed to get vault values")
		}
		return models.SecretModel{}, errors.New("vault values not found")
	}

	var values []models.ValueDTO

	for rows.Next() {
		var value models.ValueDTO

		if err := rows.Scan(&value.Key, &value.Value); err != nil {
			log.Error("failed to scan values", sl.OpErr(op, err))
			return models.SecretModel{}, err
		}

		values = append(values, value)
	}

	if rows.Err() != nil {
		log.Error("rows error", sl.OpErr(op, err))
		return models.SecretModel{}, err
	}

	res := models.ConvertDTOToSecretModel(vault, values)
	return res, nil
}

func (r *DBClient) CheckVault(ctx context.Context, log *slog.Logger, id int) error {
	const op = "db.postgresql.CheckVault"

	tx, err := r.dbClient.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", sl.OpErr(op, err))
		return errors.New("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	getVaultQuery := `
		SELECT id, name FROM vault
		WHERE id = $1;
	`

	log.Debug("get vault query", slog.String("op", op), slog.String("query", utils.QueryConvert(getVaultQuery)))

	var vault models.VaultModel
	err = tx.QueryRow(ctx, getVaultQuery, id).Scan(&vault.ID, &vault.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error("failed to get vault", sl.OpErr(op, err))
			return errors.New("vault not found")
		}
		log.Error("failed to get vault", sl.OpErr(op, err))
		return errors.New("failed to get vault")
	}
	return nil
}
