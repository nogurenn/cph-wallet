package transaction

import "github.com/google/uuid"

type Service interface {
	// CreateAccount creates a new user account.
	CreateAccount(username string) error
	// GetAccounts fetches all accounts and their respective balances.
	GetAccounts() ([]Account, error)
}

type service struct {
	db Repository
}

func NewService(db Repository) Service {
	return &service{db: db}
}

func (s *service) CreateAccount(username string) error {
	txn, err := s.db.BeginTxn()
	if err != nil {
		return err
	}

	newAccount := Account{
		Id:       uuid.New(),
		Username: username,
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
