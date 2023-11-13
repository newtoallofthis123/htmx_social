build:
	@go build -o bin/htmx_social

run: build
	@./bin/htmx_social
