package pubsub

import (
	"context"
	"net/http"
)

type ctxKey int

const (
	collKey ctxKey = iota
)

func UseCollection(coll *Collection) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), collKey, coll),
			))
		})
	}
}

func RequestCollection(r *http.Request) *Collection {
	co, ok := r.Context().Value(collKey).(*Collection)
	if !ok {
		return nil
	}
	return co
}
