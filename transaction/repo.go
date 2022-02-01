package transaction

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nogurenn/cph-wallet/dbutil"
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
)

const transactionSchemaName = "wallet"

type Repository interface {
	// BeginTxn creates a transaction object to be used by queries and commands representing a single transaction.
	BeginTxn() (dbutil.Transaction, error)
	// GetAccounts retrieves a slice of Account instances.
	GetAccounts(txn dbutil.Transaction) ([]Account, error)
	// GetAccountByUsername retrieves an Account by username.
	GetAccountByUsername(txn dbutil.Transaction, username string) (*Account, error)
	// CreateAccount creates an Account in the storage.
	CreateAccount(txn dbutil.Transaction, account Account) error
	// GetTransactionsByName retrieves all transactions with name `name` and their respective entries.
	GetTransactionsByName(txn dbutil.Transaction, name string) ([]Transaction, error)
	// LockTransactions acquires a lock for transactions to be used in conjunction with CreateTransaction.
	LockTransactions(txn dbutil.Transaction) error
	// CreateTransaction creates a Transaction in the storage, and should be used only after LockTransactions.
	CreateTransaction(txn dbutil.Transaction, transaction Transaction) error
	// CreateEntriesForTransactionId creates multiple entries under a given Transaction.
	CreateEntriesForTransactionId(txn dbutil.Transaction, transactionId uuid.UUID, entries []Entry) error
}

type postgresDb struct {
	*sqlx.DB
}

func NewPostgresDb(db *sqlx.DB) Repository {
	return &postgresDb{db}
}

func (db *postgresDb) BeginTxn() (dbutil.Transaction, error) {
	txn, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	if err := dbutil.SwitchSchema(txn, transactionSchemaName); err != nil {
		txn.Rollback()
		return nil, err
	}

	return txn, nil
}

const sqlGetAccounts = `
SELECT
	a.id,
	a.username,
	a.currency,
	COALESCE(SUM(te.credit + te.debit), 0.0) AS balance
FROM accounts a LEFT JOIN transaction_entries te ON a.id = te.account_id
GROUP BY a.id
ORDER BY a.username
`

func (db *postgresDb) GetAccounts(txn dbutil.Transaction) ([]Account, error) {
	var accounts []Account
	if err := txn.Select(&accounts, sqlGetAccounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

const sqlGetAccountByUsername = `
SELECT * FROM accounts WHERE username = $1 
`

func (db *postgresDb) GetAccountByUsername(txn dbutil.Transaction, username string) (*Account, error) {
	account := new(Account)
	if err := txn.Get(account, sqlGetAccountByUsername, username); err != nil {
		return nil, err
	}
	return account, nil
}

const sqlCreateAccount = `
INSERT INTO accounts (id, username, currency) VALUES (:id, :username, :currency)
`

func (db *postgresDb) CreateAccount(txn dbutil.Transaction, account Account) error {
	_, err := txn.NamedExec(sqlCreateAccount, account)
	return err
}

const sqlGetTransactionsByName = `
SELECT
	t.id,
	t.name,
	t.created_at,
	t.updated_at,
	te.account_id,
	te.target_account_id,
	te.name AS entry_name,
	te.credit,
	te.debit,
	a1.username,
	a2.username AS target_username
FROM transactions t 
INNER JOIN transaction_entries te ON t.id = te.transaction_id
INNER JOIN accounts a1 ON te.account_id = a1.id
LEFT OUTER JOIN accounts a2 ON te.target_account_id = a2.id
WHERE t.name = $1
GROUP BY 
	t.id,
	te.account_id,
	te.target_account_id,
	entry_name,
	te.credit,
	te.debit,
	a1.username,
	target_username
ORDER BY t.created_at DESC
`

type transactionJoinEntry struct {
	// Transaction
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	dbutil.Timestamps

	// Entry
	AccountId       uuid.UUID       `db:"account_id"`
	TargetAccountId uuid.NullUUID   `db:"target_account_id"`
	EntryName       string          `db:"entry_name"`
	Credit          decimal.Decimal `db:"credit"`
	Debit           decimal.Decimal `db:"debit"`

	// Account
	AccountName       string      `db:"username"`
	TargetAccountName null.String `db:"target_username"`
}

func (db *postgresDb) GetTransactionsByName(txn dbutil.Transaction, name string) ([]Transaction, error) {
	var transactionWithEntryRows []transactionJoinEntry
	if err := txn.Select(&transactionWithEntryRows, sqlGetTransactionsByName, name); err != nil {
		return nil, err
	}
	if transactionWithEntryRows == nil {
		return nil, nil
	}

	// since rows are sorted by this stage, we can collect all entries by transaction by folding left
	var (
		transactions           []Transaction
		previousRow            transactionJoinEntry
		lastTransactionPointer int
	)
	for i, row := range transactionWithEntryRows {
		// save the loop variable first because looping over structs in golang reuses the same memory address internally as optimization
		r := row

		// handle first row
		if i == 0 {
			transactions = append(transactions, mapRowToTransaction(r))
			previousRow = r
			continue
		}

		if r.Id == previousRow.Id {
			transactions[lastTransactionPointer].Entries = append(
				transactions[lastTransactionPointer].Entries,
				mapRowToEntry(r),
			)
		} else {
			transactions = append(transactions, mapRowToTransaction(r))
			lastTransactionPointer += 1
		}

		previousRow = r
	}

	return transactions, nil
}

// row-level lock to prevent concurrent inserts/updates
const sqlLockTransactions = `
SELECT transactions.id FROM transactions FOR UPDATE
`

func (db *postgresDb) LockTransactions(txn dbutil.Transaction) error {
	var transactionIds []uuid.UUID
	return txn.Select(&transactionIds, sqlLockTransactions)
}

const sqlCreateTransaction = `
INSERT INTO transactions (id, name) VALUES (:id, :name)
`

func (db *postgresDb) CreateTransaction(txn dbutil.Transaction, transaction Transaction) error {
	_, err := txn.NamedExec(sqlCreateTransaction, transaction)
	return err
}

const sqlCreateEntriesForTransactionId = `
INSERT INTO transaction_entries (
	id,
	transaction_id,
	account_id,
	target_account_id,
	name,
	credit,
	debit
) VALUES (
	:id,
	:transaction_id,
	:account_id,
	:target_account_id,
	:name,
	:credit,
	:debit
)
`

func (db *postgresDb) CreateEntriesForTransactionId(txn dbutil.Transaction, transactionId uuid.UUID, entries []Entry) error {
	// ensure that all entries belong to transactionId
	for _, entry := range entries {
		if entry.TransactionId != transactionId {
			return ErrTransactionEntryMismatch
		}
	}

	_, err := txn.NamedExec(sqlCreateEntriesForTransactionId, entries)
	return err
}

// --- helpers

func mapRowToTransaction(row transactionJoinEntry) Transaction {
	return Transaction{
		Id:   row.Id,
		Name: row.Name,
		Timestamps: dbutil.Timestamps{
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		},
		Entries: []Entry{mapRowToEntry(row)},
	}
}

func mapRowToEntry(row transactionJoinEntry) Entry {
	return Entry{
		AccountId:         row.AccountId,
		TargetAccountId:   row.TargetAccountId,
		Name:              row.EntryName,
		Credit:            row.Credit,
		Debit:             row.Debit,
		AccountName:       row.AccountName,
		TargetAccountName: row.TargetAccountName,
	}
}
