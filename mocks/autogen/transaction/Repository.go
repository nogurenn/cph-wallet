// Code generated by mockery v2.9.4. DO NOT EDIT.

package transaction

import (
	dbutil "github.com/nogurenn/cph-wallet/dbutil"
	mock "github.com/stretchr/testify/mock"

	transaction "github.com/nogurenn/cph-wallet/transaction"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// BeginTxn provides a mock function with given fields:
func (_m *Repository) BeginTxn() (dbutil.Transaction, error) {
	ret := _m.Called()

	var r0 dbutil.Transaction
	if rf, ok := ret.Get(0).(func() dbutil.Transaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(dbutil.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAccount provides a mock function with given fields: txn, account
func (_m *Repository) CreateAccount(txn dbutil.Transaction, account transaction.Account) error {
	ret := _m.Called(txn, account)

	var r0 error
	if rf, ok := ret.Get(0).(func(dbutil.Transaction, transaction.Account) error); ok {
		r0 = rf(txn, account)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccounts provides a mock function with given fields: txn
func (_m *Repository) GetAccounts(txn dbutil.Transaction) ([]transaction.Account, error) {
	ret := _m.Called(txn)

	var r0 []transaction.Account
	if rf, ok := ret.Get(0).(func(dbutil.Transaction) []transaction.Account); ok {
		r0 = rf(txn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]transaction.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dbutil.Transaction) error); ok {
		r1 = rf(txn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
