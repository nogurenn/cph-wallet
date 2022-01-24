//go:build integration
// +build integration

package transaction_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nogurenn/cph-wallet/dbutil"
	"github.com/nogurenn/cph-wallet/transaction"
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
