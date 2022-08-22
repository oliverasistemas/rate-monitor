package data

import "github.com/hashicorp/go-hclog"

type ExchangeDataProvider interface {
	GetExchangeRates() *ExchangeRatesList
}

type ExchangeRates struct {
	log   hclog.Logger
	rates *ExchangeRatesList
}
