package main

import (
	"coreum/conversion/data"
	protos "coreum/conversion/protos"
	"coreum/conversion/server"
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	rates, err := data.NewExchangeRates(log)
	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()

	c := server.NewCurrency(rates, log)

	protos.RegisterCurrencyServer(gs, c)

	reflection.Register(gs)

	port := 9000

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error("unable to create listener", "error", err)
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("serving gRPC server on port %d", port))
	err = gs.Serve(listen)

	if err != nil {
		log.Error("unable to serve", "error", err)
	}
}
