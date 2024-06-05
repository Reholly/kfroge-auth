-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS account (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255) NOT NULL,

    is_email_confirmed BOOLEAN NOT NULL,
    is_banned BOOLEAN NOT NULL,

    hashed_password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS code (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id uuid NOT NULL REFERENCES account(id),
    value VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';


DROP TABLE IF EXISTS account;
DROP TABLE IF EXISTS code;
-- +goose StatementEnd
