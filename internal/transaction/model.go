package transaction

import (
	"github.com/google/uuid"
	"github.com/nogurenn/cph-wallet/internal/dbutil"
	"github.com/shopspring/decimal"
)

type Account struct {
	Id                uuid.UUID `db:"id"`
	Username          string    `db:"username"`
	dbutil.Timestamps `json:"-"`
}

type Transaction struct {
	Id                uuid.UUID `db:"id"`
	Name              string    `db:"name"`
	dbutil.Timestamps `json:"-"`
}

type Entry struct {
	Id            uuid.UUID       `db:"id"`
	TransactionId uuid.UUID       `db:"transaction_id"`
	AccountId     uuid.NullUUID   `db:"account_id"`
	Name          string          `db:"name"`
	Credit        decimal.Decimal `db:"credit"`
	Debit         decimal.Decimal `db:"debit"`
}
