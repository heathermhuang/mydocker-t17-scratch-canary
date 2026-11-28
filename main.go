package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Synthetic state exists only for the lifetime of this disposable process.
type settings struct {
	mu      sync.Mutex
	Version string
	Enabled bool
}

var page = template.Must(template.New("settings").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>MyDocker acceptance fixture</title><style>body{font:18px system-ui;max-width:42rem;margin:4rem auto;padding:0 1.5rem;line-height:1.6}input,button{font:inherit;padding:.6rem}label{display:block;margin:1rem 0}.error{color:#a00}</style><main><h1>Acceptance settings</h1><p>Synthetic test data only. This fixture makes no external API requests.</p>{{if .Error}}<p role="alert" class="error">{{.Error}}</p>{{end}}<form method="post" action="/"><label><input type="checkbox" name="enabled" {{if .Enabled}}checked{{end}}> Enabled</label><label>CLI version <input name="version" value="{{.Version}}" maxlength="32" required></label><button type="submit">Save settings</button></form><p>Minimum supported version: 1.0.13. Initial version: 1.0.44.</p></main></html>`))

func supported(version string) bool {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}
	var n [3]int
	for i, part := range parts {
		if part == "" || len(part) > 5 || (len(part) > 1 && part[0] == '0') {
			return false
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return false
			}
		}
		value, err := strconv.Atoi(part)
		if err != nil || value > 65535 {
			return false
		}
		n[i] = value
	}
	return n[0] > 1 || (n[0] == 1 && (n[1] > 0 || n[2] >= 13))
}

func (s *settings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.URL.Path == "/health" {
		w.Write([]byte("ok"))
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	s.mu.Lock()
	defer s.mu.Unlock()
	message := ""
	if r.Method == http.MethodPost {
		r.Body = http.MaxBytesReader(w, r.Body, 2048)
		if r.ParseForm() != nil || !supported(r.FormValue("version")) {
			w.WriteHeader(http.StatusBadRequest)
			message = "Enter a supported version, at least 1.0.13. Settings were not changed."
		} else {
			s.Version, s.Enabled = r.FormValue("version"), r.FormValue("enabled") == "on"
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	} else if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page.Execute(w, struct {
		Version string
		Enabled bool
		Error   string
	}{s.Version, s.Enabled, message})
}

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	s := &settings{Version: "1.0.44", Enabled: true}
	server := &http.Server{Addr: ":" + port, Handler: s, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second}
	log.Fatal(server.ListenAndServe())
}
