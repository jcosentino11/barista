default:
    @just --list

build:
    go build -o bin/barista .

run: build
    ./bin/barista

test:
    go test -v -coverpkg=./barista -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

integration-test:                
    go test -v -tags=integration -coverpkg=./barista -coverprofile=coverage.integration.out ./tests
    go tool cover -html=coverage.integration.out -o coverage.integration.html

benchmark:                                        
    go test -bench=. ./... | tee benchmark.out

fmt:
    go fmt ./...

vet:
    go vet ./...

check: fmt vet test integration-test

clean:
    rm -rf bin/
    rm -f coverage.out coverage.html
    go clean
