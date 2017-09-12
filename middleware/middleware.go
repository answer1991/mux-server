package middleware

import "net/http"

func Use(h http.Handler, middleware ...func(handler http.Handler) http.Handler) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}
