//go:build integration
// +build integration

package transaction_test

import (
	"testing"

	"github.com/google/uuid"
	dbutil2 "github.com/nogurenn/cph-wallet/dbutil"
	"github.com/nogurenn/cph-wallet/transaction"
	"github.com/stretchr/testify/assert"
)

func Test_PostgresDb_GetAccounts(t *testing.T) {
	// given
	cfg := dbutil2.NewConfig()
	db, err := dbutil2.NewDb(cfg)
	assert.NoError(t, err)
	pdb := transaction.NewPostgresDb(db)

	alice := transaction.Account{Id: uuid.New(), Username: "alice456"}
	bob := transaction.Account{Id: uuid.New(), Username: "bob123"}

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
	assert.Equal(t, bob.Id, fetched[1].Id)
	assert.Equal(t, bob.Username, fetched[1].Username)
}
