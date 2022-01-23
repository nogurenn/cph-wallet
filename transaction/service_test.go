package transaction_test

import (
	"testing"

	"github.com/google/uuid"
	mockdbutil "github.com/nogurenn/cph-wallet/mocks/autogen/dbutil"
	mocktransaction "github.com/nogurenn/cph-wallet/mocks/autogen/transaction"
	"github.com/nogurenn/cph-wallet/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Service_CreateAccount_Success(t *testing.T) {
	// given
	username := "alice456"

	txn := new(mockdbutil.Transaction)
	txn.On("Commit").Return(nil)

	db := new(mocktransaction.Repository)
	db.On("BeginTxn").Return(txn, nil)
	db.On("CreateAccount",
		txn,
		mock.MatchedBy(func(account transaction.Account) bool {
			return assert.Equal(t, username, account.Username) &&
				assert.NotNil(t, account.Id) &&
				assert.Equal(t, "USD", account.Currency)
		}),
	).Return(nil)

	service := transaction.NewService(db)

	// when
	err := service.CreateAccount(username)

	// then
	assert.NoError(t, err)
}

func Test_Service_GetAccounts_Success(t *testing.T) {
	// given
	alice := transaction.Account{Id: uuid.New(), Username: "alice456", Currency: "USD"}
	bob := transaction.Account{Id: uuid.New(), Username: "bob123", Currency: "USD"}
	accounts := []transaction.Account{alice, bob}

	txn := new(mockdbutil.Transaction)
	txn.On("Rollback").Return(nil)

	db := new(mocktransaction.Repository)
	db.On("BeginTxn").Return(txn, nil)
	db.On("GetAccounts", txn).Return(accounts, nil)

	service := transaction.NewService(db)

	// when
	fetched, err := service.GetAccounts()
	assert.NoError(t, err)

	// then
	assert.Len(t, fetched, 2)

	assert.Equal(t, alice, fetched[0])
	assert.Equal(t, bob, fetched[1])
}
