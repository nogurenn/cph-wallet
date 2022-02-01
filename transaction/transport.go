package transaction

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func MakeHandler(s Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getAccountsHandler := kithttp.NewServer(
		makeGetAccountsEndpoint(s),
		decodeGetAccountsRequest,
		encodeResponse,
		opts...,
	)
	getPaymentTransactionsHandler := kithttp.NewServer(
		makeGetPaymentTransactionsEndpoint(s),
		decodeGetPaymentTransactionsRequest,
		encodeResponse,
		opts...,
	)
	sendPaymentHandler := kithttp.NewServer(
		makeSendPaymentEndpoint(s),
		decodeSendPaymentRequest,
		encodeSendPaymentResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/transaction/v1/accounts", getAccountsHandler).Methods("GET")
	r.Handle("/transaction/v1/payments", getPaymentTransactionsHandler).Methods("GET")
	r.Handle("/transaction/v1/payments", sendPaymentHandler).Methods("POST")

	return r
}

func decodeGetAccountsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return getAccountsRequest{}, nil
}

func decodeGetPaymentTransactionsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return getPaymentTransactionsRequest{}, nil
}

func decodeSendPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req sendPaymentRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeSendPaymentResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
