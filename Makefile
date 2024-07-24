build:
	@go build -o bin/json-viewer cmd/server/main.go

test:
	@go test -v ./...

run: build
	@./bin/json-viewer

