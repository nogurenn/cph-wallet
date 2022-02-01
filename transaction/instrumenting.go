package transaction

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) GetAccounts() ([]Account, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get_accounts").Add(1)
		s.requestLatency.With("method", "get_accounts").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetAccounts()
}

func (s *instrumentingService) GetPaymentTransactions() ([]Transaction, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get_payment_transactions").Add(1)
		s.requestLatency.With("method", "get_payment_transactions").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetPaymentTransactions()
}
