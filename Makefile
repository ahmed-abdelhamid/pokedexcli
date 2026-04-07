BINARY_NAME := pokedexcli

.PHONY: build test lint fmt vet clean run check

## build: Compile the binary
build:
	go build -o $(BINARY_NAME) .

## test: Run all tests
test:
	go test ./... -v -race -count=1

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## fmt: Format all Go files
fmt:
	gofmt -w .
	goimports -w .

## vet: Run go vet
vet:
	go vet ./...

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME)
	go clean

## run: Build and run the binary
run: build
	./$(BINARY_NAME)

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
