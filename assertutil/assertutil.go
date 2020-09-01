package assertutil

import (
	"net/http"
	"testing"
)

// SuccessStatus ...
func SuccessStatus(t *testing.T, r *http.Response) {
	t.Helper()
	if r.StatusCode > 299 {
		t.Fatalf("got %s want successful status", r.Status)
	}
}

// ServerErrorStatus ...
func ServerErrorStatus(t *testing.T, r *http.Response) {
	t.Helper()
	if r.StatusCode < 500 {
		t.Fatalf("got %s want server error", r.Status)
	}
}

// NotError ...
func NotError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
