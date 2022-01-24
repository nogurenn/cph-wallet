package transaction

import (
	"github.com/jmoiron/sqlx"
	"github.com/nogurenn/cph-wallet/dbutil"
)

const transactionSchemaName = "wallet"

type Repository interface {
	// BeginTxn creates a transaction object to be used by queries and commands representing a single transaction.
	BeginTxn() (dbutil.Transaction, error)
	// GetAccounts retrieves a slice of Account instances.
	GetAccounts(txn dbutil.Transaction) ([]Account, error)
	// CreateAccount creates an Account in the storage.
	CreateAccount(txn dbutil.Transaction, account Account) error
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
`

func (db *postgresDb) GetAccounts(txn dbutil.Transaction) ([]Account, error) {
	var accounts []Account
	if err := txn.Select(&accounts, sqlGetAccounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

const sqlCreateAccount = `
INSERT INTO accounts (id, username, currency) VALUES (:id, :username, :currency)
`

func (db *postgresDb) CreateAccount(txn dbutil.Transaction, account Account) error {
	_, err := txn.NamedExec(sqlCreateAccount, account)
	return err
}
