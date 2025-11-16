default:
    @just --list

build:
    go build -o bin/barista .

run:
    go run .

test:
    go test -v -coverpkg=./internal -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

fmt:
    go fmt ./...

vet:
    go vet ./...

check: fmt vet test

clean:
    rm -rf bin/
    rm -f coverage*.out coverage*.html benchmark*.out
    go clean

deps:
    go mod download
    go mod tidy
