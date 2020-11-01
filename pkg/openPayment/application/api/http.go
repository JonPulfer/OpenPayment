package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/application/service/account"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/infrastructure/eventStream"
)

type HTTPServer struct {
	stream  eventStream.EventStream
	account account.Account
}

func NewHTTPServer(stream eventStream.EventStream, account account.Account) *HTTPServer {
	return &HTTPServer{stream: stream, account: account}
}

func (hs HTTPServer) Run() error {

	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "application/json"))

	r.Route("/account", func(r chi.Router) {
		r.Post("/", hs.addAccount)
		r.Get("/{accountNumber}", hs.getAccount)
		r.Post("/{accountNumber}", hs.updateAccount)
	})

	return http.ListenAndServe(":8080", r)
}

func (hs HTTPServer) getAccount(w http.ResponseWriter, r *http.Request) {
	accountNumber := chi.URLParam(r, "accountNumber")
	accountNum, err := strconv.Atoi(accountNumber)
	if err != nil {
		http.Error(w, "invalid account number", http.StatusBadRequest)
		log.Error().Err(err).Msg("invalid account number")
		return
	}
	fetchedAccount, err := hs.account.Fetch(accountNum)
	if err != nil {
		http.Error(w, "failed to get account", http.StatusInternalServerError)
		log.Error().Err(err).Fields(map[string]interface{}{
			"accountNumber": accountNumber,
		}).Msg("failed to get account")
		return
	}
	w.Header().Add("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(&fetchedAccount)
	if err != nil {
		http.Error(w, "unable to write response data", http.StatusInternalServerError)
		log.Error().Err(err).Msg("unable to write response data")
	}
	log.Debug().Fields(map[string]interface{}{
		"accountNumber": accountNumber,
	}).Msg("processed account get request")
}

func (hs HTTPServer) addAccount(w http.ResponseWriter, r *http.Request) {
	var acc openPayment.Account
	err := json.NewDecoder(r.Body).Decode(&acc)
	if err != nil {
		http.Error(w, "failed to extract account from body", http.StatusBadRequest)
		log.Error().Err(err).Msg("failed to extract account from request body")
		return
	}
	r.Body.Close()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&acc)
	if err != nil {
		http.Error(w, "failed to create event data", http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to create event data")
		return
	}
	ev := openPayment.NewEvent(openPayment.AccountAddEvent, buf.Bytes())
	if err := hs.stream.Publish(ev); err != nil {
		http.Error(w, "failed to publish state change event", http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to publish state change event")
		return
	}
	log.Debug().Msg("processed account add request")
	requestDone(w)
}

func (hs HTTPServer) updateAccount(w http.ResponseWriter, r *http.Request) {
	accountNumber := chi.URLParam(r, "accountNumber")
	accountNum, err := strconv.Atoi(accountNumber)
	if err != nil {
		http.Error(w, "invalid account number", http.StatusBadRequest)
		log.Error().Err(err).Msg("invalid account number")
		return
	}
	fetchedAccount, err := hs.account.Fetch(accountNum)
	if err != nil {
		http.Error(w, "failed to get account", http.StatusInternalServerError)
		log.Error().Err(err).Fields(map[string]interface{}{
			"accountNumber": accountNumber,
		}).Msg("failed to get account")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&fetchedAccount)
	if err != nil {
		http.Error(w, "failed to extract account from body", http.StatusBadRequest)
		log.Error().Err(err).Msg("failed to extract account from request body")
		return
	}
	r.Body.Close()

	if fetchedAccount.Number != accountNum {
		http.Error(w, "incorrect account in payload", http.StatusBadRequest)
		log.Error().Fields(map[string]interface{}{
			"updatedAccount": fetchedAccount,
			"accountNum":     accountNum,
		}).Msg("account number in payload does not match request")
		return
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&fetchedAccount)
	if err != nil {
		http.Error(w, "failed to create update data", http.StatusInternalServerError)
		log.Error().Err(err).Fields(map[string]interface{}{
			"accountNumber":  accountNum,
			"updatedAccount": fetchedAccount,
		}).Msg("failed to create account update data")
		return
	}

	ev := openPayment.NewEvent(openPayment.AccountUpdateEvent, buf.Bytes())
	err = hs.stream.Publish(ev)
	if err != nil {
		http.Error(w, "failed to publish account update event", http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to publish account update event")
		return
	}
	log.Debug().Msg("processed account update request")
	requestDone(w)
}

func requestDone(w http.ResponseWriter) {
	fmt.Fprint(w, "accepted")
}
