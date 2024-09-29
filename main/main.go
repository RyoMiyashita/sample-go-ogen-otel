package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"

	"sample-ogen-otel/logo"
	"sample-ogen-otel/pkg/swaggerui"
)

//go:embed static
var files embed.FS

func main() {
	slog.Info("hello.")

	router, err := newRouter()
	if err != nil {
		slog.Error("failed to create server", "err", err.Error())
		os.Exit(1)
	}

	if err := http.ListenAndServe(":8080", router); err != nil {
		slog.Error("failed to listen and serve", "err", err.Error())
		os.Exit(1)
	}

	slog.Info("bye bye.")
}

func newRouter() (http.Handler, error) {
	mux := http.NewServeMux()
	api, err := newAPIRouter()
	if err != nil {
		return nil, err
	}
	mux.Handle("/api/", http.StripPrefix("/api", api))
	mux.Handle("/static/", http.FileServer(http.FS(files)))
	mux.Handle("/swagger-ui", swaggerui.HandleSwaggerUI("http://localhost:8080/static/alias_openapi.yaml"))

	return mux, nil
}

func newAPIRouter() (http.Handler, error) {
	service := NewLogoService()
	auth := NewTokenService()
	return logo.NewServer(service, auth)
}
