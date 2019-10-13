package dao_test

import (
	"testing"
	"time"

	"github.com/ankur22/ankur-curve-euro-exchange/internal/dao"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/util"
)

func TestNetworkDao(t *testing.T) {
	t.Run("perform successful exchange rate request from EUR to GBP", func(t *testing.T) {
		// given
		timeout := time.Duration(5 * time.Second)
		client := &util.ClientMock{
			Timeout: timeout,
		}
		client.Response = util.GetValidResponseForExchangeRate()
		dao := getValidNetworkDao(client)

		// when
		body, err := dao.GetExchangeRateForNow("EUR", "GBP")

		// then
		util.AssertErrorNil(t, err)
		util.AssertNotEquals(t, 0, body.Rates["GBP"])
	})

	t.Run("perform successful exchange rate request from EUR to GBP", func(t *testing.T) {
		// given
		timeout := time.Duration(5 * time.Second)
		client := &util.ClientMock{
			Timeout: timeout,
		}
		client.Response = util.GetValidResponseForExchangeRate()
		dao := getValidNetworkDao(client)

		// when
		body, err := dao.GetExchangeRateForNow("EUR", "USD")

		// then
		util.AssertErrorNil(t, err)
		util.AssertNotEquals(t, 0, body.Rates["USD"])
	})

	t.Run("ensure only status 200 is valid", func(t *testing.T) {
		// given
		timeout := time.Duration(5 * time.Second)
		client := &util.ClientMock{
			Timeout: timeout,
		}
		client.Response = util.GetInvalidResponseForExchangeRate()
		dao := getValidNetworkDao(client)

		// when
		body, err := dao.GetExchangeRateForNow("EUR", "GBP")

		// then
		util.AssertErrorNotNil(t, err)
		util.AssertTrue(t, body == nil)
	})

	t.Run("ensure only expected valid body response works", func(t *testing.T) {
		// given
		timeout := time.Duration(5 * time.Second)
		client := &util.ClientMock{
			Timeout: timeout,
		}
		client.Response = util.GetInalid200ResponseForExchangeRate()
		dao := getValidNetworkDao(client)

		// when
		body, err := dao.GetExchangeRateForNow("EUR", "GBP")

		// then
		util.AssertErrorNotNil(t, err)
		util.AssertTrue(t, body == nil)
	})

	t.Run("perform successful exchange rate request 7 days ago from EUR to GBP", func(t *testing.T) {
		// given
		timeout := time.Duration(5 * time.Second)
		client := &util.ClientMock{
			Timeout: timeout,
		}
		client.Response = util.GetValidResponseForExchangeRate()
		dao := getValidNetworkDao(client)

		// when
		body, err := dao.GetExchangeRateFromWeekAgo("EUR", "GBP", time.Now())

		// then
		util.AssertErrorNil(t, err)
		util.AssertNotEquals(t, 0, body.Rates["GBP"])
	})
}

func getValidNetworkDao(client *util.ClientMock) dao.NetworkDAO {
	return dao.CreateNewFerAPI("https://api.exchangeratesapi.io/latest", "latest", "2019-10-13", 3, client)
}
