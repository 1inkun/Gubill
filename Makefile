run:
	@swag fmt -g internal/routers/router.go
	@swag init -g internal/routers/router.go 
	@go run cmd/main/main.go
