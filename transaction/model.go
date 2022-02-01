package transaction

import (
	"github.com/google/uuid"
	"github.com/nogurenn/cph-wallet/dbutil"
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
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
	Entries           []Entry `json:"entries"`
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

	AccountName       string      `db:"username"`
	TargetAccountName null.String `db:"target_username"`
}

type Payment struct {
	Id                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	dbutil.Timestamps `json:"-"`
	Entries           []PaymentEntry `json:"entries"`
}

type PaymentEntry struct {
	Username    string          `json:"account"`
	Amount      decimal.Decimal `json:"amount"`
	ToAccount   string          `json:"to_account,omitempty"`
	FromAccount string          `json:"from_account,omitempty"`
	Direction   string          `json:"direction"`
}
