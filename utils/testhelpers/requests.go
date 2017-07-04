package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

// Request represents a request for testing
type Request struct {
	URL      string
	Code     int
	Type     string
	PostBody string
}

// PerformRequests performs a series of test requests
func PerformRequests(t testing.TB, m http.Handler, requests []Request) {
	for _, p := range requests {
		var req *http.Request
		switch p.Type {
		case "GET":
			req, _ = http.NewRequest(p.Type, p.URL, nil)
		case "POST":
			req, _ = http.NewRequest(p.Type, p.URL, strings.NewReader(p.PostBody))
		default:
			req, _ = http.NewRequest(p.Type, p.URL, nil)
		}

		rw := httptest.NewRecorder()
		m.ServeHTTP(rw, req)
		if rw.Code != p.Code {
			if _, f, l, ok := runtime.Caller(1); ok {
				t.Errorf("Expected error code %v on url %v, got error code %v (%v line: %v", p.Code, p.URL, rw.Code, f, l)
			} else {
				t.Errorf("Expected error code %v on url %v, got error code %v", p.Code, p.URL, rw.Code)
			}
		}
	}
}
