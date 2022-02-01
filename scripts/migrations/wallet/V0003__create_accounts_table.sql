CREATE TABLE accounts
(
    id         UUID PRIMARY KEY,
    username   TEXT                     NOT NULL,
    currency   TEXT                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT uq_accounts_username_currency UNIQUE (username, currency)
);

CREATE INDEX idx_accounts_username ON accounts (username);

CREATE TRIGGER set_updated_at_accounts
    BEFORE UPDATE
    ON accounts
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at_to_now();

