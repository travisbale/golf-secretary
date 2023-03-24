init:
	git config core.hooks .githooks

generate:
	go generate ./...

test:
	go test ./...

build:
	go build -o teetimer cmd/main.go 
