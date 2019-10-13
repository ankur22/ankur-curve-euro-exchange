package dao_test

import (
	"testing"

	"github.com/ankur22/ankur-curve-euro-exchange/internal/dao"
	"github.com/ankur22/ankur-curve-euro-exchange/internal/util"
)

func TestMemstore(t *testing.T) {
	t.Run("exchange data is stored", func(t *testing.T) {
		// given
		c := &util.Clock{}
		d := dao.CreateNewMemstore(c)

		// when
		d.Store("EUR", "GBP", 0.8, true)
		oneUnit, shouldExchange, dt := d.Get("EUR", "GBP")

		// then
		util.AssertEquals(t, 0.8, oneUnit)
		util.AssertTrue(t, shouldExchange)
		util.AssertNotNil(t, dt)
	})
}
