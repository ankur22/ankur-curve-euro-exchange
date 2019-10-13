package service

import (
	"time"

	"github.com/ankur22/ankur-curve-euro-exchange/internal/dao"
)

// ExchangeRateServiceResponse - Response for the exchange rate
//								 request from the service.
type ExchangeRateServiceResponse struct {
	oneUnit        float32
	shouldExchange bool
	dataDateTime   time.Time
}

// ExchangeRateService - The service that will perform the requests
//						 against the exchange rate site and decide
//						 whether it's a good idea to exchange the
//						 chosen currency.
type ExchangeRateService interface {
	PerformRequest(from, to string) (*ExchangeRateServiceResponse, error)
}

type localExchangeRateService struct {
	networkDAO *dao.NetworkDAO
	dbDAO      *dao.DatabaseDAO
}

// CreateNewExchangeRateService - Use this to create the service
//								  layer.
func CreateNewExchangeRateService(networkDAO *dao.NetworkDAO, dbDAO *dao.DatabaseDAO) *localExchangeRateService {
	service := &localExchangeRateService{networkDAO: networkDAO, dbDAO: dbDAO}
	return service
}

// PerformRequest - Get the exchange rate netween from and to.
//				    Decide if it's a good time to exchange currencies.
func (l *localExchangeRateService) PerformRequest(from, to string) (*ExchangeRateServiceResponse, error) {
	return nil, nil
}
