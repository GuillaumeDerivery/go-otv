# GO - OIDC (Open ID Connect) Token Validation

This package aims to be bound to the nginx `auth_request`(http://nginx.org/en/docs/http/ngx_http_auth_request_module.html) to validate incoming JWT Tokens in the Authorization header.
## Pre-requisite

### > GO 1.15
https://golang.org/doc/install

## Run

Run the project:
```bash
go build
AUD="YOUR_AUDIENCE" ISS="YOUR_ISSUER" JWKS_ENDPOINT="YOUR_ISSUER_JWKS_CERT_ENDPOINT" ./go-otv
```

The docker way:
```bash
docker build -t go-otv .
docker run -e AUD="YOUR_AUDIENCE" -e ISS="YOUR_ISSUER" -e JWKS_ENDPOINT="YOUR_ISSUER_JWKS_CERT_ENDPOINT" -p 8000:8000  -t go-otv
```
## Environment variables

| Key | Commentary | Default value |
|-----|------------|---------------|
| AUD | Audience of the client you asked for an token | "" |
| ISS | The OIDC issuer | "" |
| JWKS_ENDPOINT | The ISSUER endpoint | "" |

## Credits

Heavily based on the great work of lestrrat-go:  
https://github.com/lestrrat-go/jwx  
The best JW* package.
## Other

Made in 🇫🇷   
With ❤️  
And 🥐  
(And 🍷)
