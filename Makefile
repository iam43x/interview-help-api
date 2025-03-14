.PHONY: gen-jwt-rsa http-dev

gen-jwt-rsa:
	@mkdir -p certs/jwt
	@openssl genpkey -algorithm RSA -out certs/jwt/private.pem
	@openssl rsa -in certs/jwt/private.pem -pubout -out certs/jwt/public.pem

http:
	@go run cmd/http/main.go

init-db:
	@sqlite3 db/users.db < db/v1/table.sql