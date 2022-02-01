package transaction_test

import (
	"testing"

	"github.com/google/uuid"
	mockdbutil "github.com/nogurenn/cph-wallet/mocks/autogen/dbutil"
	mocktransaction "github.com/nogurenn/cph-wallet/mocks/autogen/transaction"
	"github.com/nogurenn/cph-wallet/transaction"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"
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

func Test_Service_GetPaymentTransactions_Success(t *testing.T) {
	// given
	alice := transaction.Account{
		Id:       uuid.New(),
		Username: "alice456",
	}
	bob := transaction.Account{
		Id:       uuid.New(),
		Username: "bob123",
	}

	paymentId := uuid.New()
	payment := transaction.Transaction{
		Id:   paymentId,
		Name: transaction.PaymentTransaction,
		Entries: []transaction.Entry{
			{
				AccountId: alice.Id,
				TargetAccountId: uuid.NullUUID{
					UUID:  bob.Id,
					Valid: true,
				},
				Name:              transaction.IncomingEntry,
				Credit:            decimal.NewFromFloat(100.00),
				AccountName:       alice.Username,
				TargetAccountName: null.NewString(bob.Username, true),
			}, {
				AccountId: bob.Id,
				TargetAccountId: uuid.NullUUID{
					UUID:  alice.Id,
					Valid: true,
				},
				Name:              transaction.OutgoingEntry,
				Debit:             decimal.NewFromFloat(-100.00),
				AccountName:       bob.Username,
				TargetAccountName: null.NewString(alice.Username, true),
			},
		},
	}

	txn := new(mockdbutil.Transaction)
	txn.On("Rollback").Return(nil)

	db := new(mocktransaction.Repository)
	db.On("BeginTxn").Return(txn, nil)
	db.On("GetTransactionsByName", txn, transaction.PaymentTransaction).Return([]transaction.Transaction{payment}, nil)

	service := transaction.NewService(db)

	// when
	fetched, err := service.GetPaymentTransactions()
	assert.NoError(t, err)

	// then
	assert.Len(t, fetched, 1)

	var foundIncoming, foundOutgoing int
	for _, entry := range fetched[0].Entries {
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

func Test_Service_Deposit_Success(t *testing.T) {
	// given
	alice := &transaction.Account{Id: uuid.New(), Username: "alice456", Currency: "USD"}
	amount := decimal.NewFromFloat(50.0)

	txn := new(mockdbutil.Transaction)
	txn.On("Commit").Return(nil)

	db := new(mocktransaction.Repository)
	db.On("BeginTxn").Return(txn, nil)
	db.On("GetAccountByUsername", txn, alice.Username).Return(alice, nil)
	db.On("LockTransactions", txn).Return(nil)
	db.On("CreateTransaction", txn, mock.MatchedBy(func(tr transaction.Transaction) bool {
		return assert.NotEqual(t, uuid.Nil, tr.Id) &&
			assert.Equal(t, transaction.DepositTransaction, tr.Name)
	})).Return(nil)
	db.On("CreateEntriesForTransactionId",
		txn,
		mock.MatchedBy(func(id uuid.UUID) bool {
			return assert.NotEqual(t, uuid.Nil, id)
		}),
		mock.MatchedBy(func(entries []transaction.Entry) bool {
			return assert.Len(t, entries, 1) &&
				assert.NotEqual(t, uuid.Nil, entries[0].Id) &&
				assert.NotEqual(t, uuid.Nil, entries[0].TransactionId) &&
				assert.Equal(t, alice.Id, entries[0].AccountId) &&
				assert.Equal(t, transaction.IncomingEntry, entries[0].Name) &&
				assert.True(t, entries[0].Credit.IsPositive()) &&
				assert.True(t, entries[0].Credit.Equal(amount))
		}),
	).Return(nil)

	service := transaction.NewService(db)

	// when
	err := service.Deposit(alice.Username, amount)

	// then
	assert.NoError(t, err)
}
