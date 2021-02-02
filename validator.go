package main

import (
	"net/http"
	"fmt"
	"os"
	"encoding/json"
	"log"
	"strings"
	"errors"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type ValidationConfig struct {
	Aud       string
	Iss       string
	JwksEndpoint       string
}

type GoResponse struct {
	Message    string
}

func BuildConfiguration() ValidationConfig {
	getEnv := func(key string) string{
		value := os.Getenv(key)
		if len(value) == 0 {
			log.Printf("ENV Key=%v is mandatory", key)
			return ""
		}
		return value
	}
	validationConfiguration := ValidationConfig{
		Aud: getEnv("AUD"),
		Iss: getEnv("ISS"),
		JwksEndpoint: getEnv("JWKS_ENDPOINT"),
	}
	jsonValidationConfiguration, _ := json.Marshal(validationConfiguration)
	log.Printf("Configuration loaded: %s", jsonValidationConfiguration)
	return validationConfiguration
}

func ExtractTokenFromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	xRequestId := r.Header.Get("X-Parent-Request-Id")

	log.Printf("[%s] Authorization header is: %s", xRequestId, authHeader)
	if authHeader == "" {
		log.Printf("[%s] Error: Required authorization token not found", xRequestId)
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func InitialiseJwkSet(config ValidationConfig) (*jwk.Set, error) {
	set, err := jwk.FetchHTTP(config.JwksEndpoint)
	if err != nil {
	  log.Printf("Error: Failed to parse JWK: %s", err)
	  return nil, err
	}
	log.Printf("JWKS loaded !")
	return set, nil
}

type ValidationType = func(responseWriter http.ResponseWriter, request *http.Request)

func Validate(jwks *jwk.Set, configuration ValidationConfig)  ValidationType {
	var handler = func(responseWriter http.ResponseWriter, request *http.Request) {
		xRequestId := request.Header.Get("X-Parent-Request-Id")
		// for name, values := range request.Header {
		// 	for _, value := range values {
		// 		fmt.Println(name, value)
		// 	}
		// }
		originalMethod := request.Header.Get("X-Original-Method")
		log.Printf(
			"[%s] %s %s for %s %s -- %s on %s",
			xRequestId,
			request.Method,
			request.URL.Path,
			originalMethod,
			request.Header.Get("X-Original-Url"),
			request.Header.Get("X-Forwarded-For"),
			request.Header.Get("User-Agent"),
		)

		if originalMethod == http.MethodOptions {
			log.Printf("[%s] Method OPTIONS is authorized directly", xRequestId)
			makeJsonResponse(responseWriter, http.StatusOK, "OK")
			return
		}

		// Token extraction with error management
		token, err := ExtractTokenFromAuthHeader(request)
		if err != nil {
			makeJsonResponse(responseWriter, http.StatusUnauthorized, fmt.Sprintf("Error: Extracting JWT: %v", err))
			return
		}

		// Now parse the token
		parsedToken, err := jwt.ParseString(token, jwt.WithKeySet(jwks))
		if err != nil {
			makeJsonResponse(responseWriter, http.StatusUnauthorized, fmt.Sprintf("Error: Parsing JWT: %v", err))
			return
		}
		jsonParsedToken, _ := json.Marshal(parsedToken)
		log.Printf("[%s] Decoded token extracted: %s", xRequestId, jsonParsedToken)

		err = jwt.Validate(
			parsedToken,
			jwt.WithIssuer(configuration.Iss),
			jwt.WithAudience(configuration.Aud),
		)
		if err != nil {
			makeJsonResponse(responseWriter, http.StatusUnauthorized, fmt.Sprintf("Error validating JWT: %v", err))
			return 
		}
		makeJsonResponse(responseWriter, http.StatusOK, "OK")
	}
	return handler
}

func makeJsonResponse(responseWriter http.ResponseWriter, status int, message string) {
	response := GoResponse{Message: message}
	js, err := json.Marshal(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	  }
	  if status != http.StatusOK {
	  }
	  responseWriter.WriteHeader(status)
	  responseWriter.Write(js)
}

func main() {
	configuration := BuildConfiguration()
	jwkSet, err :=  InitialiseJwkSet(configuration)
	if err != nil {
		return
	}
	http.HandleFunc("/validate", Validate(jwkSet, configuration))
	log.Fatal(http.ListenAndServe(":8000", nil))
}