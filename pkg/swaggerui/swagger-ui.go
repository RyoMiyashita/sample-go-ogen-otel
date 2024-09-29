package swaggerui

import (
	"html/template"
	"log/slog"
	"net/http"
)

func HandleSwaggerUI(srcURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t, err := template.ParseFiles("pkg/swaggerui/index.html")
		if err != nil {
			slog.ErrorContext(ctx, "failed to parse index.html", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = t.Execute(w, struct {
			SwaggerSrcURL string
		}{
			SwaggerSrcURL: srcURL,
		}); err != nil {
			slog.ErrorContext(ctx, "failed to execute template", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
