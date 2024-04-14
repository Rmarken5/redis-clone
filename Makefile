build:
	@go build -o bin/redis-clone

run: build
	@./bin/redis-clone

