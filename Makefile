.PHONY: protos

define HELP_TEXT

  Makefile commands

	make test         - Run the full test suite
	make start 		  - Build and start gRPC server
	make protos       - Compile .proto file
    make test-server  - Test grpc server (grpcurl needs to be installed)

endef

help:
	$(info $(HELP_TEXT))

test:
	go test -v --cover "./..."

protos:
	 protoc --go_out=. --go_opt=paths=source_relative \
         --go-grpc_out=. --go-grpc_opt=paths=source_relative \
         protos/currency.proto

start:
	 go build -o coreum-conversion . && ./coreum-conversion

test-server:
	grpcurl --plaintext -d '{"Base": "bnb", "Destination": "usd"}' localhost:9000 Currency/GetRate