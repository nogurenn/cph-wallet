package transaction

import (
	"time"

	"github.com/go-kit/log"
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
