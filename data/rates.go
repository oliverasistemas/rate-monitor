package data

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"strings"
	"time"
)

// NewExchangeRates create new exchange provider
func NewExchangeRates(l hclog.Logger) (*ExchangeRates, error) {
	geckoProvider := CoinGeckoProvider{}
	exchangeRates, err := geckoProvider.GetExchangeRates()
	if err != nil {
		return nil, err
	}
	rates := ExchangeRates{log: l, rates: exchangeRates}

	l.Info("fetched rates from CoinGeckoProvider")

	return &rates, nil
}

// GetRate get the exchange rate between 2 symbols
func (e *ExchangeRates) GetRate(base string, dest string) (float64, error) {
	list := *e.rates
	baseRate, ok := list[strings.ToLower(base)]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", base)
	}

	destRate, ok := list[strings.ToLower(dest)]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", dest)
	}

	return destRate.Value / baseRate.Value, nil
}

// MonitorRates checks the rates with the specified API every interval and sends a message to the
// returned channel when there are changes
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				//update rates
				geckoProvider := CoinGeckoProvider{}
				exchangeRates, err := geckoProvider.GetExchangeRates()
				if err != nil {
					e.log.Warn("an error was produced while updating the rates")
				}
				e.rates = exchangeRates
				r := (*e.rates)["usd"]
				e.log.Info(fmt.Sprintf("updating rates, new BTC_USD value: %.4f", r.Value))

				// notify updates, this will block unless there is a listener on the other end
				ret <- struct{}{}
			}
		}
	}()

	return ret
}
