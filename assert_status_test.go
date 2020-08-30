package main

import (
	"net/http/httptest"
	"testing"
)

func assertSuccessStatus(t *testing.T, r *httptest.ResponseRecorder) {
	t.Helper()
	if r.Result().StatusCode > 299 {
		t.Fatalf("got %s want successful status", r.Result().Status)
	}
}

func assertServerError(t *testing.T, r *httptest.ResponseRecorder) {
	t.Helper()
	if r.Result().StatusCode < 500 {
		t.Fatalf("got %s want server error", r.Result().Status)
	}
}
