package transaction_test

import (
	"testing"

	"github.com/google/uuid"
	mockdbutil "github.com/nogurenn/cph-wallet/mocks/autogen/dbutil"
	mocktransaction "github.com/nogurenn/cph-wallet/mocks/autogen/transaction"
	"github.com/nogurenn/cph-wallet/transaction"
	"github.com/nogurenn/cph-wallet/util"
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

	txn.AssertExpectations(t)
	db.AssertExpectations(t)
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

	txn.AssertExpectations(t)
	db.AssertExpectations(t)
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
				AccountId:         alice.Id,
				TargetAccountId:   util.NewNullUUID(bob.Id),
				Name:              transaction.IncomingEntry,
				Credit:            decimal.NewFromFloat(100.00),
				AccountName:       alice.Username,
				TargetAccountName: null.NewString(bob.Username, true),
			}, {
				AccountId:         bob.Id,
				TargetAccountId:   util.NewNullUUID(alice.Id),
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

	txn.AssertExpectations(t)
	db.AssertExpectations(t)
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

	txn.AssertExpectations(t)
	db.AssertExpectations(t)
}

func Test_Service_SendPayment_Success(t *testing.T) {
	// given
	aliceUsername := "alice456"
	bobUsername := "bob123"
	alice := &transaction.Account{Id: uuid.New(), Username: aliceUsername, Currency: "USD", Balance: decimal.NewFromFloat(200.0)}
	bob := &transaction.Account{Id: uuid.New(), Username: bobUsername, Currency: "USD", Balance: decimal.NewFromFloat(200.0)}
	amount := decimal.NewFromFloat(100.0)

	txn := new(mockdbutil.Transaction)
	txn.On("Commit").Return(nil)

	db := new(mocktransaction.Repository)
	db.On("BeginTxn").Return(txn, nil)
	db.On("GetAccountByUsername", txn, aliceUsername).Return(alice, nil).Once()
	db.On("GetAccountByUsername", txn, bobUsername).Return(bob, nil).Once()
	db.On("LockTransactions", txn).Return(nil)
	db.On("CreateTransaction", txn, mock.MatchedBy(func(tr transaction.Transaction) bool {
		return assert.NotEqual(t, uuid.Nil, tr.Id) &&
			assert.Equal(t, transaction.PaymentTransaction, tr.Name)
	})).Return(nil)
	db.On("CreateEntriesForTransactionId",
		txn,
		mock.MatchedBy(func(id uuid.UUID) bool {
			return assert.NotEqual(t, uuid.Nil, id)
		}),
		mock.MatchedBy(func(entries []transaction.Entry) bool {
			assert.Len(t, entries, 2)

			var creditEntry, debitEntry *transaction.Entry
			for _, entry := range entries {
				e := entry
				if entry.Name == transaction.OutgoingEntry {
					debitEntry = &e
				} else if entry.Name == transaction.IncomingEntry {
					creditEntry = &e
				}
			}
			assert.NotNil(t, creditEntry)
			assert.NotNil(t, debitEntry)

			assert.Equal(t, creditEntry.TransactionId, debitEntry.TransactionId)
			// credit entry
			assert.True(t, uuid.Nil != creditEntry.Id)
			assert.Equal(t, alice.Id, creditEntry.AccountId)
			assert.Equal(t, util.NewNullUUID(bob.Id), creditEntry.TargetAccountId)
			assert.Equal(t, transaction.IncomingEntry, creditEntry.Name)
			assert.True(t, creditEntry.Credit.IsPositive())
			assert.True(t, creditEntry.Debit.IsZero())
			assert.True(t, amount.Equal(creditEntry.Credit))
			// debit entry
			assert.True(t, uuid.Nil != debitEntry.Id)
			assert.Equal(t, bob.Id, debitEntry.AccountId)
			assert.Equal(t, util.NewNullUUID(alice.Id), debitEntry.TargetAccountId)
			assert.Equal(t, transaction.OutgoingEntry, debitEntry.Name)
			assert.True(t, debitEntry.Credit.IsZero())
			assert.True(t, debitEntry.Debit.IsNegative())
			assert.True(t, amount.Neg().Equal(debitEntry.Debit))

			return true
		}),
	).Return(nil)

	service := transaction.NewService(db)

	// when
	err := service.SendPayment(bobUsername, aliceUsername, amount)

	// then
	assert.NoError(t, err)

	txn.AssertExpectations(t)
	db.AssertExpectations(t)
}

func Test_Service_SendPayment_PaymentSenderReceiverIdentical(t *testing.T) {
	// given
	amount := decimal.NewFromFloat(201.0)

	db := new(mocktransaction.Repository)

	service := transaction.NewService(db)

	// when
	err := service.SendPayment("alice456 ", " alice456  ", amount)

	// then
	assert.Equal(t, transaction.ErrPaymentSenderReceiverIdentical, err)

	db.AssertExpectations(t)
}

func Test_Service_SendPayment_CreditAmountInvalid(t *testing.T) {
	// given
	amount := decimal.NewFromFloat(100.0)

	db := new(mocktransaction.Repository)

	service := transaction.NewService(db)

	// when
	err := service.SendPayment("bob123", "alice456", amount.Neg())

	// then
	assert.Equal(t, transaction.ErrCreditAmountInvalid, err)

	db.AssertExpectations(t)
}

func Test_Service_SendPayment_InsufficientBalance(t *testing.T) {
	// given
	aliceUsername := "alice456"
	bobUsername := "bob123"
	alice := &transaction.Account{Id: uuid.New(), Username: aliceUsername, Currency: "USD", Balance: decimal.NewFromFloat(200.0)}
	bob := &transaction.Account{Id: uuid.New(), Username: bobUsername, Currency: "USD", Balance: decimal.NewFromFloat(200.0)}
	amount := decimal.NewFromFloat(201.0)

	txn := new(mockdbutil.Transaction)
	txn.On("Rollback").Return(nil)

	db := new(mocktransaction.Repository)
	db.On("BeginTxn").Return(txn, nil)
	db.On("GetAccountByUsername", txn, aliceUsername).Return(alice, nil).Once()
	db.On("GetAccountByUsername", txn, bobUsername).Return(bob, nil).Once()

	service := transaction.NewService(db)

	// when
	err := service.SendPayment(bobUsername, aliceUsername, amount)

	// then
	assert.Equal(t, transaction.ErrBalanceInsufficient, err)

	txn.AssertExpectations(t)
	db.AssertExpectations(t)
}
