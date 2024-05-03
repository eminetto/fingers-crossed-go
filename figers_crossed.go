package middleware

import (
	"bufio"
	"container/ring"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

// FingersCrossed is a middleware that captures log entries and flushes them
func FingersCrossed(minLog slog.Level, triggerLog slog.Level, rng *ring.Ring, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redirect STDOUT to a buffer
		stdout := os.Stdout
		rf, wf, _ := os.Pipe()
		os.Stdout = wf
		next.ServeHTTP(w, r.WithContext(r.Context()))
		// Reset output
		wf.Close()
		os.Stdout = stdout
		flush := false
		scanner := bufio.NewScanner(rf)
		for scanner.Scan() {
			l := parseLog(scanner.Text())
			if l.Level >= minLog {
				rng.Value = l
				rng = rng.Next()
			}
			if l.Level >= triggerLog {
				flush = true
			}

		}
		if flush {
			rng.Do(func(p any) {
				if p != nil {
					fmt.Println(p)
				}
			})
			n := rng.Len()
			rng = ring.New(n)
		}
	})
}

type logEntry struct {
	Time    string     `json:"time"`
	Level   slog.Level `json:"level"`
	Message string     `json:"msg"`
}

func parseLog(raw string) logEntry {
	var l logEntry
	json.Unmarshal([]byte(raw), &l)
	return l
}
