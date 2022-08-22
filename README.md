# Exchange rates monitor
Fast and performant exchange rate queries. The monitor uses [CoinGecko API](https://www.coingecko.com/en/api) to feed gRPC server. 
To perform conversion fast we store live exchange rates in local memory and sync them periodically.

### Available scripts ###

Get some help
```
make help
```

Test the code (includes coverage)
```
make test
```

Build and start gRPC server
```
make start
```

Test grpc server (grpcurl needs to be installed)
```
make test-server
```

Compile .proto file
```
make protos
```
