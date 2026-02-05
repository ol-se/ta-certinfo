.SILENT:

mock:
	find -type d -name "mocks" -exec rm -rf "{}" + && mockery

cover:
	go test ./... -coverprofile=.coverage
	go tool cover -html=.coverage
	rm .coverage

test:
	go test ./...

lint:
	golangci-lint run ./...

build:
	go build -o ./dist/certinfo ./cmd