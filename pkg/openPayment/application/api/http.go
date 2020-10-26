package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment/infrastructure/eventStream"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/service/account"
)

type HTTPServer struct {
	stream  eventStream.EventStream
	account account.Account
}

func NewHTTPServer(stream eventStream.EventStream) *HTTPServer {
	return &HTTPServer{stream: stream}
}

func (hs HTTPServer) Run() error {

	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "application/json"))

	r.Route("/account", func(r chi.Router) {
		r.Get("/{accountNumber}", hs.getAccount)
	})

	return nil
}

func (hs HTTPServer) getAccount(w http.ResponseWriter, r *http.Request) {
	accountNumber := chi.URLParam(r, "accountNumber")
	accountNum, err := strconv.Atoi(accountNumber)
	if err != nil {
		http.Error(w, "invalid account number", http.StatusBadRequest)
		return
	}
	fetchedAccount, err := hs.account.Fetch(accountNum)
	if err != nil {
		http.Error(w, "failed to get account", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(&fetchedAccount)
	if err != nil {
		http.Error(w, "unable to write response data", http.StatusInternalServerError)
	}
}
