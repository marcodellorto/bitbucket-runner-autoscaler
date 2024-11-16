CUR_DIR := $$(pwd)

vendor:
	go mod vendor

lint:
	@golangci-lint run ./...

test:
	go test -v -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
