package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"sample-ogen-otel/logo"
)

const (
	Issuer = "logo"
	Expiry = 24 * time.Hour
)

type LogoService struct {
	logos map[string]logo.LogoDetail
	mux   sync.Mutex
}

var _ logo.Handler = (*LogoService)(nil)

func NewLogoService() *LogoService {
	return &LogoService{
		logos: make(map[string]logo.LogoDetail),
	}
}

func (l *LogoService) GetToken(ctx context.Context, req *logo.TokenRequest) (*logo.TokenResponse, error) {
	email := req.Email
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	claims := jwt.MapClaims{
		"email": email,
		"iss":   Issuer,
		"exp":   time.Now().Add(Expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(getTokenSecret()))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &logo.TokenResponse{
		Token: tokenString,
	}, nil
}

func (l *LogoService) CreateLogo(ctx context.Context, req *logo.LogoCreate) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	idUUID, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}
	id := idUUID.String()

	if _, ok := l.logos[id]; ok {
		return fmt.Errorf("logo already exists: %s", id)
	}

	now := time.Now()

	l.logos[id] = logo.LogoDetail{
		LogoId:    id,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return nil
}

func (l *LogoService) GetLogoList(ctx context.Context, params logo.GetLogoListParams) (*logo.LogoSearchResult, error) {
	l.mux.Lock()
	defer l.mux.Unlock()

	logos := make([]logo.LogoDetail, 0, len(l.logos))
	for _, l := range l.logos {
		logos = append(logos, l)
	}

	return &logo.LogoSearchResult{
		Logos:      logos,
		TotalCount: len(logos),
	}, nil
}

func (l *LogoService) NewError(ctx context.Context, err error) *logo.ErrorStatusCode {
	slog.ErrorContext(ctx, "detect api error", "err", err.Error())
	return &logo.ErrorStatusCode{
		StatusCode: 500, // TODO: error code
		Response: logo.Error{
			Code:    500,
			Message: err.Error(),
		},
	}
}
