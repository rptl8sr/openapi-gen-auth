package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	httpmiddleware "github.com/oapi-codegen/nethttp-middleware"

	"openapi-gen-auth/api"
	"openapi-gen-auth/jwt"
)

const (
	authSchemaName = "BearerAuth"
)

type contextKey string

var (
	ContextKeyUserID = contextKey("user_id")
	Scopes           = contextKey("scopes")
)

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
)

func CreateAuthMiddleware(jwtSecret string) (func(next http.Handler) http.Handler, error) {
	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := httpmiddleware.OapiRequestValidatorWithOptions(
		swagger,
		&httpmiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(jwtSecret),
			},
			SilenceServersWarning: true,
		})

	return validator, nil
}

func NewAuthenticator(jwtSecret string) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		if input.SecuritySchemeName != authSchemaName {
			return fmt.Errorf("wrong security scheme '%s' should be '%s'", input.SecuritySchemeName, authSchemaName)
		}

		token, err := GetTokenFromRequest(input.RequestValidationInput.Request)
		claims, err := jwt.ParseToken(token, jwtSecret)
		if err != nil {
			return fmt.Errorf("parsing token: %w", err)
		}

		ctx = context.WithValue(ctx, Scopes, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
		*input.RequestValidationInput.Request = *input.RequestValidationInput.Request.WithContext(ctx)

		return nil
	}
}

func GetTokenFromRequest(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", ErrInvalidAuthHeader
	}

	return strings.TrimPrefix(authHeader, prefix), nil
}
