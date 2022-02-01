CREATE TABLE transactions
(
    id         UUID PRIMARY KEY,
    name       TEXT                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_name ON transactions (name);

CREATE TRIGGER set_updated_at_transactions
    BEFORE UPDATE
    ON transactions
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at_to_now();
