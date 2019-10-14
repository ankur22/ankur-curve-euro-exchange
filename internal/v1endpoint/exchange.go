package v1endpoint

import (
	"errors"
	"fmt"

	"github.com/ankur22/ankur-curve-euro-exchange/internal/service"
	"github.com/gin-gonic/gin"
)

type v1Exchange struct {
	exchangeService service.ExchangeRateService
	validCurrencies map[string]bool
}

// CreateNewV1Exchange - Create a new endpoint for
//						 `/v1/exchange`
func CreateNewV1Exchange(exchangeService service.ExchangeRateService, validCurrencies map[string]bool) *v1Exchange {
	return &v1Exchange{exchangeService: exchangeService, validCurrencies: validCurrencies}
}

func (v *v1Exchange) PerformRequest(r *gin.Engine) {
	r.GET("/v1/exchange", func(c *gin.Context) {
		from, to, err := v.getQueryParams(c)
		if err != nil {
			v.createBadRequestResponse(c)
			return
		}

		resp, err := v.exchangeService.PerformRequest(from, to)
		if err != nil {
			v.createServerErrorResponse(c, err)
			return
		}

		v.createSuccessResponse(c, from, to, resp)
	})
}

func (v *v1Exchange) getQueryParams(c *gin.Context) (string, string, error) {
	from := c.Query("from")
	to := c.Query("to")

	_, exists := v.validCurrencies[from]
	if exists == false {
		return "", "", errors.New(fmt.Sprintf("%s is not a valid currency", from))
	}

	_, exists = v.validCurrencies[to]
	if exists == false {
		return "", "", errors.New(fmt.Sprintf("%s is not a valid currency", to))
	}

	if from == to {
		return "", "", errors.New(fmt.Sprintf("from '%s' and to are the same, they need to be different", from))
	}

	return from, to, nil
}

func (v *v1Exchange) createBadRequestResponse(c *gin.Context) {
	c.JSON(400, gin.H{
		"reason": "query params are invalid. EUR, USD and GBP are valid.",
	})
}

func (v *v1Exchange) createServerErrorResponse(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"reason": err,
	})
}

func (v *v1Exchange) createSuccessResponse(c *gin.Context, from, to string, r *service.ExchangeRateServiceResponse) {
	c.JSON(200, gin.H{
		"from":           from,
		"to":             to,
		"singleUnit":     r.OneUnit,
		"shouldExchange": r.ShouldExchange,
		"dataDateTime":   r.DataDateTime,
	})
}
