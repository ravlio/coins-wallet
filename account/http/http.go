package http

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/ravlio/wallet/account"
)

func NewHTTPHandler(endpoints account.Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/createaccount", httptransport.NewServer(endpoints.CreateAccount, DecodeCreateAccountRequest, EncodeCreateAccountResponse))
	m.Handle("/getaccount", httptransport.NewServer(endpoints.GetAccount, DecodeGetAccountRequest, EncodeGetAccountResponse))
	m.Handle("/deleteaccount", httptransport.NewServer(endpoints.DeleteAccount, DecodeDeleteAccountRequest, EncodeDeleteAccountResponse))
	m.Handle("/updateaccount", httptransport.NewServer(endpoints.UpdateAccount, DecodeUpdateAccountRequest, EncodeUpdateAccountResponse))
	m.Handle("/listaccounts", httptransport.NewServer(endpoints.ListAccounts, DecodeListAccountsRequest, EncodeListAccountsResponse))
	return m
}
func DecodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req account.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeCreateAccountResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeGetAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req account.GetAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeGetAccountResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeDeleteAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req account.DeleteAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeDeleteAccountResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeUpdateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req account.UpdateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeUpdateAccountResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeListAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req account.ListAccountsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeListAccountsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
