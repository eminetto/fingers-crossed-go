package examples

import (
	"container/ring"
	"log/slog"
	"net/http"
	"os"

	middleware "github.com/eminetto/fingers-crossed-go"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func advancedUsageExample() {
	r := chi.NewRouter()
	r.Get("/custom", func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.Log(slog.LevelInfo, "custom log level inside the handler")
		w.Write([]byte("Hello World with custom log level"))
	})
	r.Get("/conditional", func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.New()
		logger.SetLevel(logrus.WarnLevel)
		logger.Warn("conditional logging based on logrus level")
		w.Write([]byte("Hello World with conditional logging"))
	})
	r.Get("/integration", func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.Info("integration with other middleware")
		w.Write([]byte("Hello World with middleware integration"))
	})
	rng := ring.New(5000)
	fg := middleware.FingersCrossed(slog.LevelInfo, slog.LevelError, rng, r)
	http.ListenAndServe(":3001", fg)
}
