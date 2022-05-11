package middlewares

import (
	"context"
	"log"
	"net/http"
	"tribble/handlers"
	"tribble/settings"
)

func Authentication(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			handlers.HandleApiErrors(w, http.StatusUnauthorized, "")
			return
		}

		claims, err := handlers.CheckToken(token)
		if err != nil {
			log.Printf("Could not validate token: %v", err.Error())
			handlers.HandleApiErrors(w, http.StatusForbidden, "")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, settings.E, claims.Email)
		ctx = context.WithValue(ctx, settings.I, claims.ID)
		req := r.WithContext(ctx)
		handler.ServeHTTP(w, req)
	})
}
