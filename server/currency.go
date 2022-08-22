package server

import (
	"context"
	"coreum/conversion/data"
	protos "coreum/conversion/protos"
	"github.com/hashicorp/go-hclog"
	"io"
	"time"
)

// Currency gRPC server implementation
type Currency struct {
	protos.UnimplementedCurrencyServer
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

// NewCurrency creates a new Currency server
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{rates: r, log: l}

	go c.startMonitor()

	return c
}

func (c *Currency) startMonitor() {
	rates := c.rates.MonitorRates(10 * time.Second)

	c.log.Info("starting exchange rates monitor")

	for range rates {
		c.log.Info("updated rates.")

		for k, v := range c.subscriptions {

			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase(), rr.GetDestination())
				if err != nil {
					c.log.Error("Unable to get update rate", "base", rr.GetBase(), "destination", rr.GetDestination())
				}

				err = k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r})
				if err != nil {
					c.log.Error("Unable to send updated rate", "base", rr.GetBase(), "destination", rr.GetDestination())
				}
			}
		}

	}
}

func (c *Currency) GetRate(_ context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase(), rr.GetDestination())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	for {
		rr, err := src.Recv()

		if err == io.EOF {
			c.log.Info("Client has closed connection")
			delete(c.subscriptions, src)
			break
		}

		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			delete(c.subscriptions, src)
			return err
		}

		c.log.Info("Handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}

	return nil
}

//func BatchConversion(context.Context, *protos.BatchConversionRequest) (*protos.BatchConversionResponse, error) {
//	return nil, nil
//}
//
//func RateList(context.Context, *protos.RateListRequest) (*protos.RateListResponse, error) {
//	return nil, nil
//}
