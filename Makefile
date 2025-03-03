.PHONY: gen-jwt-rsa

gen-jwt-rsa:
	@mkdir -p certs/jwt
	@openssl genpkey -algorithm RSA -out certs/jwt/private.pem
	@openssl rsa -in private.pem -pubout -out certs/jwt/public.pem