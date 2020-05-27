.PHONY: test build install

test:
	go test -v -cover ./...

pre-commit:
	go mod tidy
	go mod vendor
	go vet ./...
	go fmt ./...

install:
	go install -mod vendor

build:
	go build ./...
