package transaction

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service interface {
	// CreateAccount creates a new user account.
	CreateAccount(username string) error
	// GetAccounts fetches all accounts and their respective balances.
	GetAccounts() ([]Account, error)
	// GetPaymentTransactions fetches all transactions with name PaymentTransaction.
	GetPaymentTransactions() ([]Transaction, error)
	// Deposit records a deposit transaction for the given username, if the account exists.
	Deposit(username string, amount decimal.Decimal) error
}

type service struct {
	db Repository
}

func NewService(db Repository) Service {
	return &service{db: db}
}

const (
	defaultAccountCurrency = "USD"

	// list of valid transaction names
	PaymentTransaction = "payment"
	DepositTransaction = "deposit"

	// list of valid entry names
	IncomingEntry = "incoming"
	OutgoingEntry = "outgoing"
)

func (s *service) CreateAccount(username string) error {
	txn, err := s.db.BeginTxn()
	if err != nil {
		return err
	}

	newAccount := Account{
		Id:       uuid.New(),
		Username: username,
		Currency: defaultAccountCurrency,
	}

	if err = s.db.CreateAccount(txn, newAccount); err != nil {
		txn.Rollback()
		return err
	}

	return txn.Commit()
}

func (s *service) GetAccounts() ([]Account, error) {
	txn, err := s.db.BeginTxn()
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	return s.db.GetAccounts(txn)
}

func (s *service) Deposit(username string, amount decimal.Decimal) error {
	txn, err := s.db.BeginTxn()
	if err != nil {
		return err
	}

	if amount.IsNegative() || amount.IsZero() {
		return ErrCreditAmountInvalid
	}

	account, err := s.db.GetAccountByUsername(txn, username)
	if err != nil {
		return err
	}

	err = s.db.LockTransactions(txn)
	if err != nil {
		txn.Rollback()
		return err
	}

	depositId := uuid.New()
	err = s.db.CreateTransaction(txn, Transaction{
		Id:   depositId,
		Name: DepositTransaction,
	})
	if err != nil {
		txn.Rollback()
		return err
	}

	err = s.db.CreateEntriesForTransactionId(txn, depositId, []Entry{{
		Id:            uuid.New(),
		TransactionId: depositId,
		AccountId:     account.Id,
		Name:          IncomingEntry,
		Credit:        amount,
	}})
	if err != nil {
		txn.Rollback()
		return err
	}

	return txn.Commit()
}

func (s *service) GetPaymentTransactions() ([]Transaction, error) {
	txn, err := s.db.BeginTxn()
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	return s.db.GetTransactionsByName(txn, PaymentTransaction)
}
