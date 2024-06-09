generate:
	templ generate

build: generate
	go build

run: generate
	go run cmd/app/main.go serve
