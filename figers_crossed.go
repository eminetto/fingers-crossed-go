package middleware

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func FingersCrossed(minLog slog.Level, triggerLog slog.Level, logs []LogEntry, next http.Handler) http.Handler {
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
				logs = append(logs, l)
			}
			if l.Level >= triggerLog {
				flush = true
			}

		}
		if flush {
			for _, l := range logs {
				fmt.Println(l)
			}
			logs = nil
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
