.PHONY: lint clean deps mocks test docs run build stop

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

docs:
	swag init -g cmd/server/main.go

run:
	docker-compose up

build:
	docker-compose up --build

stop:
	docker-compose down
