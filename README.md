# GO - OIDC (Open ID Connect) Token Validation

This package aims to be bound to the nginx `auth_request`(http://nginx.org/en/docs/http/ngx_http_auth_request_module.html) to validate incoming JWT Tokens in the Authorization header.

It's validating the token's signature thanks to the exposed OIDC jwks endpoint as well as the expiry date, audience and issuer.

This package has been built to integrate a kubernetes environment and to work with the default nginx ingress.

## Kubernetes integration

By putting this line in the `Ingress` you can it protect with a mandatory Authorization Bearer token.

```yml
nginx.ingress.kubernetes.io/auth-url: http://GO-OTV-SERVICE.NAMESPACE.svc.cluster.local/validate
```

If you have public routes to handle, just create a new `Ingress` for the same host without the previous line.  

### To go further with Ingress & tracing

By putting this line in the `Ingress` you can pass the auth module the generated x-request-id from the parent client request.

```yml
nginx.ingress.kubernetes.io/auth-snippet: |
    proxy_set_header X-Parent-Request-Id $req_id;
```

It will prefix the request logs.
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
docker run -e AUD="YOUR_AUDIENCE" -e ISS="YOUR_ISSUER" -e JWKS_ENDPOINT="YOUR_ISSUER_JWKS_CERTS_ENDPOINT" -p 8000:8000  -t go-otv
```
## Environment variables

| Key | Commentary | Default value |
|-----|------------|---------------|
| AUD | Token's Audience  | "" |
| ISS | Token's Issuer | "" |
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
