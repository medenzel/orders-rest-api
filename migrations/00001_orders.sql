-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    state SMALLINT NOT NULL,
    create_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
