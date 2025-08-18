default:
    @just --list

build:
    go build -o bin/barista .

run: build
    ./bin/barista

test:
    go test ./...

test-coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

integration-test:                                  
    go test -tags=integration ./...

fmt:
    go fmt ./...

vet:
    go vet ./...

check: fmt vet test

clean:
    rm -rf bin/
    rm -f coverage.out coverage.html
    go clean

deps:
    go mod download
    go mod tidy
