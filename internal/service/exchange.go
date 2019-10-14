package service

import (
	"fmt"
	"time"

	"github.com/ankur22/ankur-curve-euro-exchange/internal/dao"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/util"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
)

// ExchangeRateServiceResponse - Response for the exchange rate
//								 request from the service.
type ExchangeRateServiceResponse struct {
	OneUnit        float32
	ShouldExchange bool
	DataDateTime   time.Time
}

type grResponse struct {
	res *dao.ExchangeRateResponse
	err error
}

// ExchangeRateService - The service that will perform the requests
//						 against the exchange rate site and decide
//						 whether it's a good idea to exchange the
//						 chosen currency.
type ExchangeRateService interface {
	PerformRequest(from, to string) (*ExchangeRateServiceResponse, error)
}

type localExchangeRateService struct {
	networkDAO        dao.NetworkDAO
	dbDAO             dao.DatabaseDAO
	dataValidDuration time.Duration
	clock             *util.Clock
	timeout           time.Duration
	sem               *semaphore.Weighted
}

// CreateNewExchangeRateService - Use this to create the service
//								  layer.
func CreateNewExchangeRateService(networkDAO dao.NetworkDAO,
	dbDAO dao.DatabaseDAO,
	dataValidDuration time.Duration,
	clock *util.Clock,
	timeout time.Duration) *localExchangeRateService {

	return &localExchangeRateService{networkDAO: networkDAO,
		dbDAO:             dbDAO,
		dataValidDuration: dataValidDuration,
		clock:             clock,
		timeout:           timeout,
		sem:               semaphore.NewWeighted(1)}
}

// PerformRequest - Get the exchange rate netween from and to.
//				    Decide if it's a good time to exchange currencies.
func (l *localExchangeRateService) PerformRequest(from, to string) (*ExchangeRateServiceResponse, error) {
	oneUnit, shouldExchange, dataDateTime := l.dbDAO.Get(from, to)
	if l.hasStoredValueExpired(dataDateTime) {
		// Only allow one thread to perform
		// the network request and save to the db
		if !l.sem.TryAcquire(1) {
			// If no data available then return error
			if dataDateTime.IsZero() {
				return nil, errors.New("Timed out waiting for another thread to complete network request")
			}
			// Use expired data
			return &ExchangeRateServiceResponse{DataDateTime: dataDateTime, OneUnit: oneUnit, ShouldExchange: shouldExchange}, nil
		}
		defer l.sem.Release(1)
		oneUnit, shouldExchange, dataDateTime, err := l.getAndStoreNewValues(from, to)
		if err != nil {
			return nil, err
		}
		return &ExchangeRateServiceResponse{DataDateTime: dataDateTime, OneUnit: oneUnit, ShouldExchange: shouldExchange}, nil
	}
	return &ExchangeRateServiceResponse{DataDateTime: dataDateTime, OneUnit: oneUnit, ShouldExchange: shouldExchange}, nil
}

func (l *localExchangeRateService) hasStoredValueExpired(dataDateTime time.Time) bool {
	now := l.clock.Now()
	diff := now.Sub(dataDateTime)
	return diff > l.dataValidDuration
}

func (l *localExchangeRateService) getAndStoreNewValues(from, to string) (float32, bool, time.Time, error) {
	chan1 := make(chan grResponse)
	chan2 := make(chan grResponse)

	go l.performLatestRequest(from, to, chan1)
	go l.performWeekAgoRequest(from, to, chan2)

	var latest *dao.ExchangeRateResponse = nil
	var weekOld *dao.ExchangeRateResponse = nil

	select {
	case res := <-chan1:
		if res.err != nil {
			return 0, false, time.Time{}, res.err
		}
		latest = res.res
	case <-time.After(l.timeout):
		return 0, false, time.Time{}, errors.New("Timeout occured while waiting for response from network layer")
	}

	select {
	case res := <-chan2:
		if res.err != nil {
			return 0, false, time.Time{}, res.err
		}
		weekOld = res.res
	case <-time.After(l.timeout):
		return 0, false, time.Time{}, errors.New("Timeout occured while waiting for response from network layer")
	}

	shouldExchange := weekOld.Rates[to] > latest.Rates[to]
	dataDateTime := l.clock.Now()
	l.dbDAO.Store(from, to, latest.Rates[to], shouldExchange, dataDateTime)

	return latest.Rates[to], shouldExchange, dataDateTime, nil
}

func (l *localExchangeRateService) performLatestRequest(from, to string, c chan grResponse) {
	respLatest, err := l.networkDAO.GetExchangeRateForNow(from, to)
	l.extractResponse(respLatest, c, err, to)
}

func (l *localExchangeRateService) performWeekAgoRequest(from, to string, c chan grResponse) {
	weekAgo := l.getDateFromWeekAgo()
	respLatest, err := l.networkDAO.GetExchangeRateFromPast(from, to, weekAgo)
	l.extractResponse(respLatest, c, err, to)
}

func (l *localExchangeRateService) extractResponse(resp *dao.ExchangeRateResponse, c chan grResponse, err error, to string) {
	if err != nil {
		c <- grResponse{nil, errors.Wrap(err, "Cannot get week old requested data")}
		return
	}

	_, exists := resp.Rates[to]
	if !exists {
		c <- grResponse{nil, errors.New(fmt.Sprintf("Response doesn't contain conversion value to '%s'", to))}
		return
	}

	c <- grResponse{resp, err}
}

func (l *localExchangeRateService) getDateFromWeekAgo() time.Time {
	now := l.clock.Now()
	return now.AddDate(0, 0, -7)
}
