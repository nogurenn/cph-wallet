-- double-entry bookkeeping
CREATE TABLE transaction_entries
(
    id                UUID PRIMARY KEY,
    transaction_id    UUID                     NOT NULL,
    -- account that owns the entry
    account_id        UUID                     NOT NULL,
    -- sender/receiver from the perspective of account_id
    target_account_id UUID,
    name              TEXT                     NOT NULL,
    credit            DECIMAL(32, 8)           NOT NULL DEFAULT 0.0 CHECK (credit >= 0.0),
    debit             DECIMAL(32, 8)           NOT NULL DEFAULT 0.0 CHECK (debit <= 0.0),
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_transaction_entries_transaction_id
        FOREIGN KEY (transaction_id) REFERENCES transactions (id)
            ON UPDATE RESTRICT
            ON DELETE RESTRICT,

    CONSTRAINT fk_transaction_entries_account_id
        FOREIGN KEY (account_id) REFERENCES accounts (id)
            ON UPDATE RESTRICT
            ON DELETE RESTRICT,

    CONSTRAINT fk_transaction_entries_target_account_id
        FOREIGN KEY (target_account_id) REFERENCES accounts (id)
            ON UPDATE RESTRICT
            ON DELETE RESTRICT
);

CREATE INDEX idx_transaction_entries_name ON transaction_entries (name);

CREATE TRIGGER set_updated_at_transaction_entries
    BEFORE UPDATE
    ON transaction_entries
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at_to_now();
