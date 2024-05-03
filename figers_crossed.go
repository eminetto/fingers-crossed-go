package middleware

import (
	"bufio"
	"bytes"
	"container/ring"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func FingersCrossed(minLog slog.Level, triggerLog slog.Level, rng *ring.Ring, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}

		// Redirect STDOUT to a buffer
		stdout := os.Stdout
		rf, wf, err := os.Pipe()
		if err != nil {
			panic(err) //@todo
		}
		os.Stdout = wf
		next.ServeHTTP(w, r.WithContext(r.Context()))
		// Reset output
		wf.Close()
		os.Stdout = stdout
		flush := false
		scanner := bufio.NewScanner(rf)
		for scanner.Scan() {
			buf.WriteString(scanner.Text())
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

type LogEntry struct {
	Time    string     `json:"time"`
	Level   slog.Level `json:"level"`
	Message string     `json:"msg"`
}

func parseLog(raw string) LogEntry {
	var l LogEntry
	json.Unmarshal([]byte(raw), &l)
	return l
}
