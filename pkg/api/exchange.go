package api

// ExchangeResponse - Reponse model of /v1/exchange
type ExchangeResponse struct {
	CurrencyFrom string
	CurrencyTo   string
	OneUnit      float32
	Buy          bool
	DataDateTime string
}
