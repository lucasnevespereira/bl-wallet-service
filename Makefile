.PHONY: lint build docs run serve test clean deps mocks test-coverage

lint:
	golangci-lint run --exclude-use-default=true

clean:
		rm -rf mocks

deps: clean
		go get -d -v ./...

mocks: clean deps
		mockery --all

test: mocks
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total | awk '{print $$3}'

docs:
	swag init -g cmd/server/main.go

run:
	docker-compose up

build:
	docker-compose up --build

stop:
	docker-compose down