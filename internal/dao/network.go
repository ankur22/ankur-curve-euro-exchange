package dao

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// ExchangeRateResponse - Model response from exchange
type ExchangeRateResponse struct {
	Base  string
	Date  string
	Rates map[string]float32
}

// NetworkDAO - interface to get exchange data from over the
//				network (e.g. HTTP, FTP etc.)
type NetworkDAO interface {
	GetExchangeRateForNow(from, to string) (*ExchangeRateResponse, error)
	GetExchangeRateFromPast(from, to string, date time.Time) (*ExchangeRateResponse, error)
}

// HTTPClient - Http client interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// FerAPI - Gets exchange rates from https://exchangeratesapi.io/
type ferAPI struct {
	url        string
	client     HTTPClient
	latest     string
	layoutISO  string
	retryCount int
}

// CreateNewFerAPI - Create new ferAPI to get exchange rates
//					 from https://exchangeratesapi.io/
func CreateNewFerAPI(url, latest, layoutISO string, retryCount int, client HTTPClient) *ferAPI {
	return &ferAPI{url: url, client: client, layoutISO: layoutISO, retryCount: retryCount, latest: latest}
}

// GetExchangeRateForNow - Get the exchange rate for {from} to {to}
func (f *ferAPI) GetExchangeRateForNow(from, to string) (*ExchangeRateResponse, error) {
	return f.getRequest(from, to, f.latest)
}

// GetExchangeRateFromPast - Get the exchange rate for {from} to {to} from n days back
func (f *ferAPI) GetExchangeRateFromPast(from, to string, date time.Time) (*ExchangeRateResponse, error) {
	test := date.Format(f.layoutISO)
	return f.getRequest(from, to, test)
}

func (f *ferAPI) getRequest(from, to, endpoint string) (*ExchangeRateResponse, error) {
	resp, err := f.performGetRequest(from, to, endpoint)
	if err != nil {
		return nil, err
	}

	data, err := f.unmarshallResponseBody(resp, from, to)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f *ferAPI) performGetRequest(from, to, endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", f.url+"/"+endpoint, nil)
	q := req.URL.Query()
	q.Add("base", from)
	q.Add("symbols", to)
	req.URL.RawQuery = q.Encode()

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Cannot get exchange rate from '%s' to '%s'", from, to))
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Received %d when getting exchange rate from '%s' to '%s'", resp.StatusCode, from, to))
	}

	return resp, nil
}

func (f *ferAPI) unmarshallResponseBody(resp *http.Response, from, to string) (*ExchangeRateResponse, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Cannot read body for exchange rate request from '%s' to '%s'", from, to))
	}

	data := ExchangeRateResponse{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Cannot unmarshall body for exchange rate request from '%s' to '%s'", from, to))
	}

	return &data, nil
}
