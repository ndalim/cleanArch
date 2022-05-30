package middleware

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"redditapp/pkg/auth/usecase"
	"redditapp/tools"
	"strings"
	"time"
)

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Auth(auth usecase.AuthUsecase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		headerAuth, ok := r.Header["Authorization"]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(headerAuth[0], " ")
		if len(headerParts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		usr, err := auth.Check(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), tools.CurrentUserKey, usr)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AccessLog(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		next.ServeHTTP(w, r)

		logger.Infow("http request",
			"url", r.URL.Path,
			"method", r.Method,
			"state", "end",
			"time", time.Since(start))

	})
}

func WrapHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("access-control-allow-origin", "*")
		w.Header().Set("cache-control", "public, max-age=0, must-revalidate")
		w.Header().Set("content-type", "application/json; charset=utf-8")

		next.ServeHTTP(w, r)
	}
}
