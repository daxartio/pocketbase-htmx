generate:
	templ generate

run: generate
	go run cmd/app/main.go serve
