package transaction

import (
	"time"

	"github.com/go-kit/log"
	"github.com/shopspring/decimal"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) GetAccounts() ([]Account, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get_accounts",
			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.GetAccounts()
}

func (s *loggingService) GetPaymentTransactions() ([]Transaction, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get_payment_transactions",
			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.GetPaymentTransactions()
}

func (s *loggingService) SendPayment(username string, targetUsername string, amount decimal.Decimal) error {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send_payment",
			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.SendPayment(username, targetUsername, amount)
}
