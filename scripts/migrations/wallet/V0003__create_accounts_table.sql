CREATE TABLE accounts
(
    id         UUID PRIMARY KEY,
    username   TEXT                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT uq_accounts_username UNIQUE (username)
);

CREATE TRIGGER set_updated_at_accounts
    BEFORE UPDATE
    ON accounts
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at_to_now();

