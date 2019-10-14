package api

// ExchangeResponse - Reponse model of /v1/exchange
type ExchangeResponse struct {
	From           string
	To             string
	SingleUnit     float32
	ShouldExchange bool
	DataDateTime   string
}

// ExchangeErrorResponse - Error reponse model of /v1/exchange
type ExchangeErrorResponse struct {
	Reason string
}
