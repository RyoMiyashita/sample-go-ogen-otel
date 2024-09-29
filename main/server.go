package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"

	"sample-ogen-otel/logo"
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
