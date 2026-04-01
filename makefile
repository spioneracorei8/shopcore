dev:
	clear
	go run cmd/api/main.go

test: 
	go test -v -cover ./...