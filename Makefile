include .env
export

.PHONY: gen-jwt-rsa http init-db docker

gen-jwt-rsa:
	@mkdir -p certs/jwt
	@openssl genpkey -algorithm RSA -out ${JWT_PRIVATE_KEY}
	@openssl rsa -in ${JWT_PRIVATE_KEY} -pubout -out ${JWT_PUBLIC_KEY}

http:
	@go run cmd/http/main.go

init-db:
	@sqlite3 ${DB_PATH} < db/v1/table.sql

docker:
	@docker build . -t api:latest