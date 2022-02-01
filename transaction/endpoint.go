package transaction

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/nogurenn/cph-wallet/dbutil"
)

type getAccountsRequest struct{}

type getAccountsResponse struct {
	Accounts []Account `json:"accounts"`
	Err      error     `json:"error"`
}

func (r getAccountsResponse) error() error { return r.Err }

func makeGetAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		_ = request.(getAccountsRequest)
		accounts, err := s.GetAccounts()
		if accounts == nil {
			accounts = []Account{} // serialize nil slice such that `"accounts": []` instead of null
		}
		return getAccountsResponse{Accounts: accounts, Err: err}, nil
	}
}

type getPaymentTransactionsRequest struct{}

type getPaymentTransactionsResponse struct {
	Payments []Payment `json:"payments"`
	Err      error     `json:"error"`
}

func (r getPaymentTransactionsResponse) error() error { return r.Err }

func makeGetPaymentTransactionsEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		_ = request.(getPaymentTransactionsRequest)
		paymentTransactions, err := s.GetPaymentTransactions()

		payments := []Payment{}
		for _, pt := range paymentTransactions {
			payments = append(payments, mapTransactionToPayment(pt))
		}

		return getPaymentTransactionsResponse{Payments: payments, Err: err}, nil
	}
}

// --- helpers

func mapTransactionToPayment(transaction Transaction) Payment {
	return Payment{
		Id:   transaction.Id,
		Name: transaction.Name,
		Timestamps: dbutil.Timestamps{
			CreatedAt: transaction.CreatedAt,
			UpdatedAt: transaction.UpdatedAt,
		},
		Entries: mapEntriesToPaymentEntries(transaction.Entries),
	}
}

func mapEntriesToPaymentEntries(entries []Entry) []PaymentEntry {
	paymentEntries := []PaymentEntry{}
	for _, entry := range entries {
		paymentEntry := PaymentEntry{
			Username:  entry.AccountName,
			Direction: entry.Name,
		}

		if paymentEntry.Direction == IncomingEntry {
			paymentEntry.Amount = entry.Credit
			paymentEntry.FromAccount = entry.TargetAccountName
		} else {
			paymentEntry.Amount = entry.Debit.Abs()
			paymentEntry.ToAccount = entry.TargetAccountName
		}

		paymentEntries = append(paymentEntries, paymentEntry)
	}

	return paymentEntries
}
