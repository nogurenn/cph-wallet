package transaction

import (
	"context"

	"github.com/go-kit/kit/endpoint"
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
