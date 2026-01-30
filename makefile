dev:
	go run cmd/main.go
gen-docs:
	swag init -g cmd/main.go