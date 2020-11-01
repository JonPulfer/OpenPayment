package main

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment/application/api"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/application/eventHandler"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/application/service/account"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/infrastructure/eventStream"
)

func main() {

	zerolog.LevelFieldName = "severity"
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	ims := eventStream.NewInMemoryStream()
	handler := eventHandler.NewHandler(ims)
	simpleAccount := account.NewSimpleAccount()
	handler.Subscribe(simpleAccount)
	var handlerWg sync.WaitGroup
	handlerWg.Add(1)
	go func() {
		err := handler.Handle(&handlerWg)
		if err != nil {
			panic(err)
		}
	}()

	svc := api.NewHTTPServer(ims, simpleAccount)
	err := svc.Run()
	if err != nil {
		log.Error().Err(err).Msg("web service failed")
	}
}
