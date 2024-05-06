-- +goose Up
-- +goose StatementBegin
CREATE TABLE short(
    id SERIAL,
    correlation_id uuid NOT NULL DEFAULT gen_random_uuid (),
    original_url text,
    short_url  text,
    date_create TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
