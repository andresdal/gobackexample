build:
	@go build -o bin/gobackexample/cmd/main.go

test:
	@go test -v ./...

run: 
	@./bin/gobackexample