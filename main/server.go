package main

import (
	"context"
	"sync"

	"sample-ogen-otel/logo"
)

type LogoService struct {
	logos map[int64]logo.LogoDetail
	mux   sync.Mutex
}

var _ logo.Handler = (*LogoService)(nil)

func NewLogoService() *LogoService {
	return &LogoService{
		logos: make(map[int64]logo.LogoDetail),
	}
}

func (l *LogoService) CreateLogo(ctx context.Context, req *logo.LogoCreate) error {
	// TODO implement me
	panic("implement me")
}

func (l *LogoService) GetLogoList(ctx context.Context, params logo.GetLogoListParams) (*logo.LogoSearchResult, error) {
	// TODO implement me
	panic("implement me")
}

func (l *LogoService) NewError(ctx context.Context, err error) *logo.ErrorStatusCode {
	// TODO implement me
	panic("implement me")
}
