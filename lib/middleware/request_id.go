package middleware

import (
	"context"
	"net/http"

	uuidGen "github.com/google/uuid"
)

// RequestIdHandler is a middleware that adds a unique request ID to the request context.
func RequestIdHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get(string(RequestIdKey))
		if id == "" {
			id = uuid()
		}
		w.Header().Set(string(RequestIdKey), id)
		ctx = context.WithValue(ctx, RequestIdKey, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func uuid() string {
	return uuidGen.New().String()
}
