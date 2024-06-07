generate: 
	templ generate

build: generate
	go build

run: generate
	go run main.go serve
