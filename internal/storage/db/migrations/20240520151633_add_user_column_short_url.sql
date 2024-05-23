-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE short ADD user_id INTEGER DEFAULT 0;

CREATE TABLE users
(
    id          SERIAL PRIMARY KEY NOT NULL,
    name        text   NOT NULL,
    pass        text   NOT NULL,
    date_create timestamp default now() ,
    is_active   bool      DEFAULT true
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
