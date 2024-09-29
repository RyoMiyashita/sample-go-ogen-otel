// Code generated by ogen, DO NOT EDIT.

package logo

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// CreateLogo implements createLogo operation.
//
// Create a new logo.
//
// POST /logos
func (UnimplementedHandler) CreateLogo(ctx context.Context, req *LogoCreate) error {
	return ht.ErrNotImplemented
}

// GetLogoList implements getLogoList operation.
//
// Get collections of logos.
//
// GET /logos
func (UnimplementedHandler) GetLogoList(ctx context.Context, params GetLogoListParams) (r *LogoSearchResult, _ error) {
	return r, ht.ErrNotImplemented
}

// NewError creates *ErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r *ErrorStatusCode) {
	r = new(ErrorStatusCode)
	return r
}
