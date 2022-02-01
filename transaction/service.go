package transaction

import (
	"strings"

	"github.com/google/uuid"
	"github.com/nogurenn/cph-wallet/util"
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
	// SendPayment records a fund transfer from one account to another.
	SendPayment(fromUsername string, toUsername string, amount decimal.Decimal) error
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

func (s *service) GetPaymentTransactions() ([]Transaction, error) {
	txn, err := s.db.BeginTxn()
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	return s.db.GetTransactionsByName(txn, PaymentTransaction)
}

func (s *service) Deposit(username string, amount decimal.Decimal) error {
	txn, err := s.db.BeginTxn()
	if err != nil {
		return err
	}

	if amount.IsNegative() || amount.IsZero() {
		txn.Rollback()
		return ErrCreditAmountInvalid
	}

	account, err := s.db.GetAccountByUsername(txn, username)
	if err != nil {
		txn.Rollback()
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

	err = s.db.CreateEntriesForTransactionId(txn, depositId, []Entry{
		newCreditEntry(depositId, account.Id, util.NewNullUUID(uuid.Nil), amount),
	})
	if err != nil {
		txn.Rollback()
		return err
	}

	return txn.Commit()
}

func (s *service) SendPayment(fromUsername string, toUsername string, amount decimal.Decimal) error {
	if amount.IsNegative() || amount.IsZero() {
		return ErrCreditAmountInvalid
	}

	sanitizedFromUsername := strings.TrimSpace(fromUsername)
	sanitizedToUsername := strings.TrimSpace(toUsername)
	if sanitizedFromUsername == sanitizedToUsername {
		return ErrPaymentSenderReceiverIdentical
	}

	txn, err := s.db.BeginTxn()
	if err != nil {
		return err
	}

	sender, err := s.db.GetAccountByUsername(txn, sanitizedFromUsername)
	if err != nil {
		txn.Rollback()
		return err
	}

	receiver, err := s.db.GetAccountByUsername(txn, sanitizedToUsername)
	if err != nil {
		txn.Rollback()
		return err
	}

	if sender.Balance.LessThan(amount) {
		txn.Rollback()
		return ErrBalanceInsufficient
	}

	err = s.db.LockTransactions(txn)
	if err != nil {
		txn.Rollback()
		return err
	}

	paymentId := uuid.New()
	err = s.db.CreateTransaction(txn, Transaction{
		Id:   paymentId,
		Name: PaymentTransaction,
	})
	if err != nil {
		txn.Rollback()
		return err
	}

	err = s.db.CreateEntriesForTransactionId(txn, paymentId, []Entry{
		newDebitEntry(paymentId, sender.Id, util.NewNullUUID(receiver.Id), amount),
		newCreditEntry(paymentId, receiver.Id, util.NewNullUUID(sender.Id), amount),
	})
	if err != nil {
		txn.Rollback()
		return err
	}

	return txn.Commit()
}

// --- helpers

func newCreditEntry(transactionId uuid.UUID, accountId uuid.UUID, targetAccountId uuid.NullUUID, amount decimal.Decimal) Entry {
	return Entry{
		Id:              uuid.New(),
		TransactionId:   transactionId,
		AccountId:       accountId,
		TargetAccountId: targetAccountId,
		Name:            IncomingEntry,
		Credit:          amount.Abs(),
	}
}

func newDebitEntry(transactionId uuid.UUID, accountId uuid.UUID, targetAccountId uuid.NullUUID, amount decimal.Decimal) Entry {
	sanitizedAmount := amount
	if amount.IsPositive() {
		sanitizedAmount = sanitizedAmount.Neg()
	}

	return Entry{
		Id:              uuid.New(),
		TransactionId:   transactionId,
		AccountId:       accountId,
		TargetAccountId: targetAccountId,
		Name:            OutgoingEntry,
		Debit:           sanitizedAmount,
	}
}
