package middleware

import (
	"context"
	"net/http"
)

func DefaultTemplate(template string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), Key("template"), template)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
