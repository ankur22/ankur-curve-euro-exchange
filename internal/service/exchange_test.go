package service_test

import (
	"testing"
	"time"

	"github.com/ankur22/ankur-curve-euro-exchange/internal/dao"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/service"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/util"
	"github.com/pkg/errors"
)

func TestExchangeRateService(t *testing.T) {
	t.Run("ensure new values retrieved if none in DB", func(t *testing.T) {
		// given
		clock := util.CreateNewClock()
		dbDao := dao.CreateNewMemstore()
		networkDao := givenValidNetworkDao()
		service := service.CreateNewExchangeRateService(&networkDao, dbDao, time.Duration(time.Second), clock, time.Duration(time.Second*5))

		// when
		resp, err := service.PerformRequest("EUR", "GBP")

		// then
		util.AssertErrorNil(t, err)
		util.AssertFalse(t, resp == nil)
		util.AssertFalse(t, resp.ShouldExchange)
	})

	t.Run("ensure cached values retrieved if valid in DB", func(t *testing.T) {
		// given
		clock := util.CreateNewClock()
		dbDao := dao.CreateNewMemstore()
		networkDao := givenValidNetworkDao()
		service := service.CreateNewExchangeRateService(&networkDao, dbDao, time.Duration(time.Second), clock, time.Duration(time.Second*5))
		resp1, _ := service.PerformRequest("EUR", "GBP")
		networkDao.resetFlags()
		time.Sleep(time.Millisecond * 500)

		// when
		resp2, err := service.PerformRequest("EUR", "GBP")

		// then
		util.AssertErrorNil(t, err)
		util.AssertFalse(t, resp2 == nil)
		util.AssertFalse(t, resp2.ShouldExchange)
		util.AssertTrue(t, resp1.DataDateTime == resp2.DataDateTime)
		util.AssertFalse(t, networkDao.latestCalled)
		util.AssertFalse(t, networkDao.weekOldCalled)
	})

	t.Run("ensure new values retrieved when cached values are invalid in DB", func(t *testing.T) {
		// given
		clock := util.CreateNewClock()
		dbDao := dao.CreateNewMemstore()
		networkDao := givenValidNetworkDao()
		service := service.CreateNewExchangeRateService(&networkDao, dbDao, time.Duration(time.Millisecond*200), clock, time.Duration(time.Second*5))
		resp1, _ := service.PerformRequest("EUR", "GBP")
		networkDao.resetFlags()
		time.Sleep(time.Millisecond * 500)

		// when
		resp2, err := service.PerformRequest("EUR", "GBP")

		// then
		util.AssertErrorNil(t, err)
		util.AssertFalse(t, resp2 == nil)
		util.AssertFalse(t, resp2.ShouldExchange)
		util.AssertFalse(t, resp1.DataDateTime == resp2.DataDateTime)
		util.AssertTrue(t, networkDao.latestCalled)
		util.AssertTrue(t, networkDao.weekOldCalled)
	})

	t.Run("ensure error returned if response from network is missing latest data", func(t *testing.T) {
		// given
		clock := util.CreateNewClock()
		dbDao := dao.CreateNewMemstore()
		networkDao := givenInvalidLatestNetworkDao()
		service := service.CreateNewExchangeRateService(&networkDao, dbDao, time.Duration(time.Millisecond*200), clock, time.Duration(time.Second*5))

		// when
		_, err := service.PerformRequest("EUR", "GBP")

		// then
		util.AssertErrorNotNil(t, err)
	})

	t.Run("ensure error returned if response from network is missing week old data", func(t *testing.T) {
		// given
		clock := util.CreateNewClock()
		dbDao := dao.CreateNewMemstore()
		networkDao := givenInvalidWeekOldNetworkDao()
		service := service.CreateNewExchangeRateService(&networkDao, dbDao, time.Duration(time.Millisecond*200), clock, time.Duration(time.Second*5))

		// when
		_, err := service.PerformRequest("EUR", "GBP")

		// then
		util.AssertErrorNotNil(t, err)
	})

	t.Run("ensure error returned if network call for latest data fails", func(t *testing.T) {
		// given
		clock := util.CreateNewClock()
		dbDao := dao.CreateNewMemstore()
		networkDao := givenNetworkServiceDownDuringLatestDataRequest()
		service := service.CreateNewExchangeRateService(&networkDao, dbDao, time.Duration(time.Millisecond*200), clock, time.Duration(time.Second*5))

		// when
		_, err := service.PerformRequest("EUR", "GBP")

		// then
		util.AssertErrorNotNil(t, err)
	})

	t.Run("ensure error returned if network call for week old data fails", func(t *testing.T) {
		// given
		clock := util.CreateNewClock()
		dbDao := dao.CreateNewMemstore()
		networkDao := givenNetworkServiceDownDuringWeekAgoDataRequest()
		service := service.CreateNewExchangeRateService(&networkDao, dbDao, time.Duration(time.Millisecond*200), clock, time.Duration(time.Second*5))

		// when
		_, err := service.PerformRequest("EUR", "GBP")

		// then
		util.AssertErrorNotNil(t, err)
	})
}

type mockNetworkDAO struct {
	latest        *dao.ExchangeRateResponse
	weekOld       *dao.ExchangeRateResponse
	latestCalled  bool
	weekOldCalled bool
}

func (m *mockNetworkDAO) resetFlags() {
	m.latestCalled = false
	m.weekOldCalled = false
}

func (m *mockNetworkDAO) GetExchangeRateForNow(from, to string) (*dao.ExchangeRateResponse, error) {
	m.latestCalled = true
	if m.latest != nil {
		return m.latest, nil
	} else {
		return nil, errors.New("Network service down")
	}
}

func (m *mockNetworkDAO) GetExchangeRateFromPast(from, to string, date time.Time) (*dao.ExchangeRateResponse, error) {
	m.weekOldCalled = true
	if m.weekOld != nil {
		return m.weekOld, nil
	} else {
		return nil, errors.New("Network service down")
	}
}

func givenValidNetworkDao() mockNetworkDAO {
	latestRates := map[string]float32{"GBP": 0.9}
	weekAgoRates := map[string]float32{"GBP": 0.8}
	latest := dao.ExchangeRateResponse{Base: "EUR", Date: "2019-10-14", Rates: latestRates}
	weekOld := dao.ExchangeRateResponse{Base: "EUR", Date: "2019-10-07", Rates: weekAgoRates}
	return mockNetworkDAO{&latest, &weekOld, false, false}
}

func givenInvalidLatestNetworkDao() mockNetworkDAO {
	latestRates := map[string]float32{"USD": 0.9}
	weekAgoRates := map[string]float32{"GBP": 0.8}
	latest := dao.ExchangeRateResponse{Base: "EUR", Date: "2019-10-14", Rates: latestRates}
	weekOld := dao.ExchangeRateResponse{Base: "EUR", Date: "2019-10-07", Rates: weekAgoRates}
	return mockNetworkDAO{&latest, &weekOld, false, false}
}

func givenInvalidWeekOldNetworkDao() mockNetworkDAO {
	latestRates := map[string]float32{"GBP": 0.9}
	weekAgoRates := map[string]float32{"USD": 0.8}
	latest := dao.ExchangeRateResponse{Base: "EUR", Date: "2019-10-14", Rates: latestRates}
	weekOld := dao.ExchangeRateResponse{Base: "EUR", Date: "2019-10-07", Rates: weekAgoRates}
	return mockNetworkDAO{&latest, &weekOld, false, false}
}

func givenNetworkServiceDownDuringLatestDataRequest() mockNetworkDAO {
	weekAgoRates := map[string]float32{"USD": 0.8}
	weekOld := dao.ExchangeRateResponse{Base: "EUR", Date: "2019-10-07", Rates: weekAgoRates}
	return mockNetworkDAO{nil, &weekOld, false, false}
}

func givenNetworkServiceDownDuringWeekAgoDataRequest() mockNetworkDAO {
	latestRates := map[string]float32{"GBP": 0.9}
	latest := dao.ExchangeRateResponse{Base: "EUR", Date: "2019-10-14", Rates: latestRates}
	return mockNetworkDAO{&latest, nil, false, false}
}
