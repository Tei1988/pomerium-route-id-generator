.PHONY: build test clean

BINARY_NAME=pomerium-route-id-gen

build:
	go build -o $(BINARY_NAME) ./cmd/pomerium-route-id-gen

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)
