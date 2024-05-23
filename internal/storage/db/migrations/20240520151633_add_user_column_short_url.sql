-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE short ADD user_id INTEGER;

CREATE TABLE users
(
    id          SERIAL PRIMARY KEY NOT NULL,
    name        text   NOT NULL,
    pass        text   NOT NULL,
    date_create timestamp default now() ,
    is_active   bool      DEFAULT true
);

ALTER TABLE short ADD CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
