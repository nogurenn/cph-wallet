//go:build integration
// +build integration

package transaction_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nogurenn/cph-wallet/dbutil"
	"github.com/nogurenn/cph-wallet/transaction"
	"github.com/nogurenn/cph-wallet/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_PostgresDb_GetAccounts(t *testing.T) {
	// given
	cfg := dbutil.NewConfig()
	db, err := dbutil.NewDb(cfg)
	assert.NoError(t, err)
	pdb := transaction.NewPostgresDb(db)

	alice := transaction.Account{Id: uuid.New(), Username: "alice456", Currency: "USD"}
	bob := transaction.Account{Id: uuid.New(), Username: "bob123", Currency: "USD"}

	// when
	txn, err := pdb.BeginTxn()
	assert.NoError(t, err)

	err = pdb.CreateAccount(txn, bob)
	assert.NoError(t, err)
	err = pdb.CreateAccount(txn, alice)
	assert.NoError(t, err)

	fetched, err := pdb.GetAccounts(txn)
	assert.NoError(t, err)

	txn.Rollback()

	// then
	assert.Len(t, fetched, 2)

	// check for correctness and sort order (ASC)
	assert.Equal(t, alice.Id, fetched[0].Id)
	assert.Equal(t, alice.Username, fetched[0].Username)
	assert.True(t, alice.Balance.IsZero())
	assert.Equal(t, bob.Id, fetched[1].Id)
	assert.Equal(t, bob.Username, fetched[1].Username)
	assert.True(t, bob.Balance.IsZero())
}

func Test_PostgresDb_GetAccountByUsername(t *testing.T) {
	// given
	cfg := dbutil.NewConfig()
	db, err := dbutil.NewDb(cfg)
	assert.NoError(t, err)
	pdb := transaction.NewPostgresDb(db)

	alice := transaction.Account{Id: uuid.New(), Username: "alice456", Currency: "USD"}
	bob := transaction.Account{Id: uuid.New(), Username: "bob123", Currency: "USD"}

	// when
	txn, err := pdb.BeginTxn()
	assert.NoError(t, err)

	err = pdb.CreateAccount(txn, bob)
	assert.NoError(t, err)
	err = pdb.CreateAccount(txn, alice)
	assert.NoError(t, err)

	fetched, err := pdb.GetAccountByUsername(txn, alice.Username)
	assert.NoError(t, err)

	txn.Rollback()

	// then
	assert.Equal(t, alice.Id, fetched.Id)
	assert.Equal(t, alice.Username, fetched.Username)
}

func Test_PostgresDb_CreateAndGetTransactionsByName(t *testing.T) {
	// given
	cfg := dbutil.NewConfig()
	db, err := dbutil.NewDb(cfg)
	assert.NoError(t, err)
	pdb := transaction.NewPostgresDb(db)

	alice := transaction.Account{Id: uuid.New(), Username: "alice456", Currency: "USD"}
	bob := transaction.Account{Id: uuid.New(), Username: "bob123", Currency: "USD"}

	// initial balances
	aliceInitialBalanceId := uuid.New()
	aliceInitialBalance := transaction.Transaction{
		Id:   aliceInitialBalanceId,
		Name: transaction.DepositTransaction,
		Entries: []transaction.Entry{
			{
				Id:            uuid.New(),
				TransactionId: aliceInitialBalanceId,
				AccountId:     alice.Id,
				Name:          transaction.IncomingEntry,
				Credit:        decimal.NewFromFloat(200.00),
			},
		},
	}
	bobInitialBalanceId := uuid.New()
	bobInitialBalance := transaction.Transaction{
		Id:   bobInitialBalanceId,
		Name: transaction.DepositTransaction,
		Entries: []transaction.Entry{
			{
				Id:            uuid.New(),
				TransactionId: bobInitialBalanceId,
				AccountId:     bob.Id,
				Name:          transaction.IncomingEntry,
				Credit:        decimal.NewFromFloat(200.00),
			},
		},
	}

	// payment transaction
	paymentId := uuid.New()
	fromBob := transaction.Entry{
		Id:              uuid.New(),
		TransactionId:   paymentId,
		AccountId:       bob.Id,
		TargetAccountId: util.NewNullUUID(alice.Id),
		Name:            transaction.OutgoingEntry,
		Debit:           decimal.NewFromFloat(-100.00),
	}
	toAlice := transaction.Entry{
		Id:              uuid.New(),
		TransactionId:   paymentId,
		AccountId:       alice.Id,
		TargetAccountId: util.NewNullUUID(bob.Id),
		Name:            transaction.IncomingEntry,
		Credit:          decimal.NewFromFloat(100.00),
	}
	payment := transaction.Transaction{
		Id:      paymentId,
		Name:    transaction.PaymentTransaction,
		Entries: []transaction.Entry{fromBob, toAlice},
	}

	// when
	txn, err := pdb.BeginTxn()
	assert.NoError(t, err)

	err = pdb.CreateAccount(txn, bob)
	assert.NoError(t, err)
	err = pdb.CreateAccount(txn, alice)
	assert.NoError(t, err)

	err = pdb.LockTransactions(txn)
	assert.NoError(t, err)

	err = pdb.CreateTransaction(txn, aliceInitialBalance)
	assert.NoError(t, err)
	err = pdb.CreateEntriesForTransactionId(txn, aliceInitialBalance.Id, aliceInitialBalance.Entries)
	assert.NoError(t, err)

	err = pdb.CreateTransaction(txn, bobInitialBalance)
	assert.NoError(t, err)
	err = pdb.CreateEntriesForTransactionId(txn, bobInitialBalance.Id, bobInitialBalance.Entries)
	assert.NoError(t, err)

	err = pdb.CreateTransaction(txn, payment)
	assert.NoError(t, err)
	err = pdb.CreateEntriesForTransactionId(txn, payment.Id, payment.Entries)
	assert.NoError(t, err)

	accounts, err := pdb.GetAccounts(txn)
	assert.NoError(t, err)

	payments, err := pdb.GetTransactionsByName(txn, transaction.PaymentTransaction)
	assert.NoError(t, err)

	txn.Rollback()

	// then
	assert.Len(t, accounts, 2)
	assert.Equal(t, alice.Username, accounts[0].Username)
	assert.True(t, accounts[0].Balance.Equal(decimal.NewFromFloat(300.00)))
	assert.Equal(t, bob.Username, accounts[1].Username)
	assert.True(t, accounts[1].Balance.Equal(decimal.NewFromFloat(100.00)))

	assert.Len(t, payments, 1)
	assert.Equal(t, payment.Id, payments[0].Id)
	assert.Len(t, payments[0].Entries, 2)

	var foundIncoming, foundOutgoing int
	for _, entry := range payments[0].Entries {
		if entry.Name == transaction.IncomingEntry {
			foundIncoming += 1
			assert.True(t, entry.Credit.IsPositive())
		} else if entry.Name == transaction.OutgoingEntry {
			foundOutgoing += 1
			assert.True(t, entry.Debit.IsNegative())
		}
		assert.True(t, entry.TargetAccountId.Valid)
		assert.True(t, entry.TargetAccountName.Valid)
		assert.NotEqual(t, entry.TargetAccountName.String, entry.AccountName)
	}
	assert.Equal(t, 1, foundIncoming)
	assert.Equal(t, 1, foundOutgoing)
}
