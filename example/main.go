package main

import (
	"container/ring"
	"log/slog"
	"net/http"
	"os"

	middleware "github.com/eminetto/fingers-crossed-go"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.Info("info inside the handler")
		w.Write([]byte("Hello World with info"))
	})
	r.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.Error("error inside the handler")
		w.Write([]byte("Hello World with error"))
	})
	r.Get("/debug", func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.Debug("debug inside the handler")
		w.Write([]byte("Hello World with debug"))
	})
	r.Get("/warn", func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.Warn("warn inside the handler")
		w.Write([]byte("Hello World with warn"))
	})
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic inside the handler")
	})

	rng := ring.New(5000)
	fg := middleware.FingersCrossed(slog.LevelInfo, slog.LevelError, rng, r)
	http.ListenAndServe(":3000", fg)
}
