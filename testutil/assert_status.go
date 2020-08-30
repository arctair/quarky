package testutil

import (
	"net/http/httptest"
	"testing"
)

// AssertSuccessStatus ...
func AssertSuccessStatus(t *testing.T, r *httptest.ResponseRecorder) {
	t.Helper()
	if r.Result().StatusCode > 299 {
		t.Fatalf("got %s want successful status", r.Result().Status)
	}
}

// AssertServerError ...
func AssertServerError(t *testing.T, r *httptest.ResponseRecorder) {
	t.Helper()
	if r.Result().StatusCode < 500 {
		t.Fatalf("got %s want server error", r.Result().Status)
	}
}
