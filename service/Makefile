.PHONY: lint
lint:
	# runs linter
	golangci-lint run --timeout=1m --build-tags integration ./...
.PHONY: test
test:
	# runs all the tests
	go test -v ./...

.PHONY: test-coveralls
test-coveralls:
	# runs all the tests with coverage and sends to coveralls
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: test-race
test-race:
	go test -v ./...


.PHONY: build
build:
	# builds the binary from service/cmd/main.go
	go build -o $(BINARY_NAME) -v ./cmd/main.go

.PHONY: run
run:
	# runs the binary
	go run ./cmd/main.go
