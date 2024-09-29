package main

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"sample-ogen-otel/logo"
)

const (
	secretEnv = "TOKEN_SECRET"
)

type TokenService struct {
	secret string
}

var _ logo.SecurityHandler = (*TokenService)(nil)

func NewTokenService() *TokenService {
	return &TokenService{secret: getTokenSecret()}
}

func getTokenSecret() string {
	secret, ok := os.LookupEnv(secretEnv)
	if !ok {
		panic("should set TOKEN_SECRET")
	}

	return secret
}

func (ts *TokenService) HandleBearerAuth(ctx context.Context, operationName string, t logo.BearerAuth) (context.Context, error) {
	token, err := jwt.Parse(
		t.Token, func(_ *jwt.Token) (interface{}, error) {
			return []byte(ts.secret), nil
		},
	)
	if err != nil {
		return ctx, err
	}
	if !token.Valid {
		return ctx, fmt.Errorf("invalid token")
	}

	return ctx, nil
}
