package v1endpoint_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/ankur22/ankur-curve-euro-exchange/internal/service"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/util"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/v1endpoint"
	"github.com/ankur22/ankur-curve-euro-exchange/pkg/api"
)

func TestExchangeEndpoint(t *testing.T) {
	t.Run("ensure 200 response and valid body when successful operations", func(t *testing.T) {
		// given
		eService := givenValidExchangeService()
		endpoint := v1endpoint.CreateNewV1Exchange(&eService, givenValidCuirrenciesList())
		server := service.CreateNewServer()
		server.Register("GET /v1/exchange", endpoint)
		go server.Start()
		time.Sleep(time.Millisecond * 500)

		// when
		body := performGetRequest(t, "EUR", "GBP", 200)
		data := unmarshalSuccess(t, "EUR", "GBP", body)

		time.Sleep(time.Millisecond * 500)

		// then
		server.Stop(time.Duration(time.Second))
		util.AssertTrue(t, data.SingleUnit == 0.8)
	})

	t.Run("ensure 400 response when bad queries passed", func(t *testing.T) {
		// given
		eService := givenValidExchangeService()
		endpoint := v1endpoint.CreateNewV1Exchange(&eService, givenValidCuirrenciesList())
		server := service.CreateNewServer()
		server.Register("GET /v1/exchange", endpoint)
		go server.Start()
		time.Sleep(time.Millisecond * 500)

		// when
		body := performGetRequest(t, "FOO", "BAR", 400)
		data := unmarshalFail(t, "FOO", "BAR", body)

		time.Sleep(time.Millisecond * 500)

		// then
		server.Stop(time.Duration(time.Second))
		util.AssertTrue(t, data.Reason == "query params are invalid. EUR, USD and GBP are valid.")
	})

	t.Run("ensure 400 response when from and to are the same", func(t *testing.T) {
		// given
		eService := givenValidExchangeService()
		endpoint := v1endpoint.CreateNewV1Exchange(&eService, givenValidCuirrenciesList())
		server := service.CreateNewServer()
		server.Register("GET /v1/exchange", endpoint)
		go server.Start()
		time.Sleep(time.Millisecond * 500)

		// when
		body := performGetRequest(t, "EUR", "EUR", 400)
		data := unmarshalFail(t, "EUR", "EUR", body)

		time.Sleep(time.Millisecond * 500)

		// then
		server.Stop(time.Duration(time.Second))
		util.AssertTrue(t, data.Reason == "query params are invalid. EUR, USD and GBP are valid.")
	})
}

type mockExchangeService struct {
	resp *service.ExchangeRateServiceResponse
	err  error
}

func (m *mockExchangeService) PerformRequest(from, to string) (*service.ExchangeRateServiceResponse, error) {
	return m.resp, m.err
}

func givenValidCuirrenciesList() map[string]bool {
	return map[string]bool{"EUR": true, "USD": true, "GBP": true}
}

func givenValidExchangeService() mockExchangeService {
	resp := &service.ExchangeRateServiceResponse{OneUnit: 0.8, ShouldExchange: true, DataDateTime: time.Now()}
	return mockExchangeService{resp, nil}
}

func performGetRequest(t *testing.T, from, to string, expectStatus int) []byte {
	t.Helper()

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", "http://0.0.0.0:8080/v1/exchange", nil)
	q := req.URL.Query()
	q.Add("from", from)
	q.Add("to", to)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(fmt.Sprintf("Cannot get exchange rate from '%s' to '%s'", from, to))
	}

	if resp.StatusCode != expectStatus {
		t.Fatal(fmt.Sprintf("Received %d when getting exchange rate from '%s' to '%s'", resp.StatusCode, from, to))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(fmt.Sprintf("Cannot read body for exchange rate request from '%s' to '%s'", from, to))
	}

	return body
}

func unmarshalSuccess(t *testing.T, from, to string, body []byte) *api.ExchangeResponse {
	t.Helper()

	data := api.ExchangeResponse{}

	err := json.Unmarshal(body, &data)
	if err != nil {
		t.Fatal(fmt.Sprintf("Cannot unmarshall body for exchange rate request from '%s' to '%s'", from, to))
	}

	return &data
}

func unmarshalFail(t *testing.T, from, to string, body []byte) *api.ExchangeErrorResponse {
	t.Helper()

	data := api.ExchangeErrorResponse{}

	err := json.Unmarshal(body, &data)
	if err != nil {
		t.Fatal(fmt.Sprintf("Cannot unmarshall body for exchange rate request from '%s' to '%s'", from, to))
	}

	return &data
}
