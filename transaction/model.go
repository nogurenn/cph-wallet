package transaction

import (
	"github.com/google/uuid"
	"github.com/nogurenn/cph-wallet/dbutil"
	"github.com/shopspring/decimal"
)

type Account struct {
	Id                uuid.UUID       `db:"id" json:"-"`
	Username          string          `db:"username" json:"id"`
	Balance           decimal.Decimal `db:"balance" json:"balance"` // decimal.Decimal marshals to string to prevent silent precision loss
	Currency          string          `db:"currency" json:"currency"`
	dbutil.Timestamps `json:"-"`
}

type Transaction struct {
	Id                uuid.UUID `db:"id"`
	Name              string    `db:"name"`
	dbutil.Timestamps `json:"-"`
}

type Entry struct {
	Id                uuid.UUID       `db:"id"`
	TransactionId     uuid.UUID       `db:"transaction_id"`
	AccountId         uuid.UUID       `db:"account_id"`        // owner of the entry
	TargetAccountId   uuid.NullUUID   `db:"target_account_id"` // sender/receiver from the perspective of AccountId
	Name              string          `db:"name"`
	Credit            decimal.Decimal `db:"credit"`
	Debit             decimal.Decimal `db:"debit"`
	dbutil.Timestamps `json:"-"`
}
