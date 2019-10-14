package main

import (
	"net/http"
	"time"

	"github.com/ankur22/ankur-curve-euro-exchange/internal/dao"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/service"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/util"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/v1endpoint"
)

func main() {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	ferDao := dao.CreateNewFerAPI("https://api.exchangeratesapi.io", "latest", "2006-01-02", client)
	dbDao := dao.CreateNewMemstore()
	clock := util.CreateNewClock()
	exchangeService := service.CreateNewExchangeRateService(ferDao, dbDao, time.Duration(time.Second*5), clock, time.Duration(time.Second*5))
	validCurrencies := map[string]bool{"EUR": true, "USD": true, "GBP": true}
	exchangeEndpoint := v1endpoint.CreateNewV1Exchange(exchangeService, validCurrencies)
	server := service.CreateNewServer()
	server.Register("GET /v1/exchange", exchangeEndpoint)
	server.Start()
}
