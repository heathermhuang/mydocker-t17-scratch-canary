package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVersionValidationAndPersistence(t *testing.T) {
	s := &settings{Version: "1.0.44", Enabled: true}
	for _, value := range []string{"1.0.12", "0.9.99", "1.0.013", "bad", "1.0.13-beta"} {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("version="+value))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		s.ServeHTTP(w, r)
		if w.Code != 400 || s.Version != "1.0.44" || !s.Enabled {
			t.Fatalf("invalid value changed settings: %q", value)
		}
	}
	for _, value := range []string{"1.0.13", "1.0.44"} {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("version="+value+"&enabled=on"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		s.ServeHTTP(w, r)
		if w.Code != 303 {
			t.Fatalf("save returned %d", w.Code)
		}
		w = httptest.NewRecorder()
		s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
		if !strings.Contains(w.Body.String(), `value="`+value+`"`) {
			t.Fatal("saved setting did not survive reload")
		}
	}
}
