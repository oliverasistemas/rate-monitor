package data

import (
	gecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

type CoinGeckoProvider struct{}

type ExchangeRatesList map[string]ExchangeRatesItemStruct

type ExchangeRatesItemStruct struct {
	Name  string  `json:"name"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

// GetExchangeRates Get BTC-to-Currency exchange rates
func (provider *CoinGeckoProvider) GetExchangeRates() (*ExchangeRatesList, error) {
	cg := newClient()

	rate, err := cg.ExchangeRates()
	if err != nil {
		return nil, err
	}

	exchangeRates := adaptExchangeRates(rate)
	return &exchangeRates, nil
}

// We encapsulate specific provider (CoinGecko) details,
// so that if we later decide to switch data providers we don't have to modify all the consumers.
func adaptExchangeRates(rate *types.ExchangeRatesItem) ExchangeRatesList {
	item := ExchangeRatesList{}

	for k, v := range *rate {
		item[k] = ExchangeRatesItemStruct{
			Name:  v.Name,
			Unit:  v.Unit,
			Value: v.Value,
			Type:  v.Type,
		}
	}
	return item
}

// Create a new CoinGecko API client.
func newClient() *gecko.Client {
	return gecko.NewClient(nil)
}
