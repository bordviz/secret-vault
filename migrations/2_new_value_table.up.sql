CREATE TABLE IF NOT EXISTS value(
    id SERIAL PRIMARY KEY NOT NULL,
    vault_id INTEGER NOT NULL REFERENCES vault(id),
    key VARCHAR NOT NULL,
    value VARCHAR NOT NULL
);
