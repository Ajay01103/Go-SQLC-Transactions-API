-- +goose Up

-- +goose StatementBegin
CREATE TABLE users (
	id              text        PRIMARY KEY,
	name            text        NOT NULL,
	email           text        NOT NULL UNIQUE,
	password        text        NOT NULL,
	profile_picture text,
	created_at      timestamptz NOT NULL DEFAULT now(),
	updated_at      timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
