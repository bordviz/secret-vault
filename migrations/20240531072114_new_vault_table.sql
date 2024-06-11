-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vault(
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL
);
CREATE INDEX IF NOT EXISTS inx_id ON vault(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vault;
DROP INDEX IF EXISTS inx_id;
-- +goose StatementEnd
