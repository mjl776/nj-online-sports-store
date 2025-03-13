build:
	@go build -o bin/nj-online-sports-store

run: build
	@./bin/nj-online-sports-store

test:
	@go test -v ./...Make