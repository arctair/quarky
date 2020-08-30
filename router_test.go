package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// StubBuild ...
type StubBuild struct {
	sha1    string
	version string
}

func (b *StubBuild) getSha1() string {
	return b.sha1
}

func (b *StubBuild) getVersion() string {
	return b.version
}

type StubDeployments struct {
	stubError string
}

func (d StubDeployments) Create() (string, error) {
	var err error
	if d.stubError != "" {
		err = errors.New(d.stubError)
	} else {
		err = nil
	}
	return "6ed4fdb9-2934-406f-a2bc-0e7cd8f301ae", err
}

func (d StubDeployments) Delete() (string, error) {
	var err error
	if d.stubError != "" {
		err = errors.New(d.stubError)
	} else {
		err = nil
	}
	return "1ed4fdb9-2934-406f-a2bc-0e7cd8f301ae", err
}

type MockLogger struct {
	errors []error
}

func NewMockLogger() MockLogger {
	return MockLogger{
		[]error{},
	}
}

func (l *MockLogger) error(err error) {
	l.errors = append(l.errors, err)
}

func (l *MockLogger) assertErrors(t *testing.T, errors []error) {
	t.Helper()
	if !reflect.DeepEqual(l.errors, errors) {
		t.Errorf("got %v want %v", l.errors, errors)
	}
}

func TestController(t *testing.T) {
	t.Run("POST / creates deployment", func(t *testing.T) {
		controller := &Router{
			nil,
			StubDeployments{""},
			nil,
		}

		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		var got map[string]string
		if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}

		want := map[string]string{
			"id": "6ed4fdb9-2934-406f-a2bc-0e7cd8f301ae",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("POST / gets deployments error", func(t *testing.T) {
		mockLogger := NewMockLogger()
		controller := &Router{
			nil,
			StubDeployments{"Stub error"},
			&mockLogger,
		}

		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		if response.Result().StatusCode < 500 {
			t.Errorf("got %s want server error", response.Result().Status)
		}

		mockLogger.assertErrors(t, []error{errors.New("Stub error")})
	})

	t.Run("DELETE / deletes deployment", func(t *testing.T) {
		controller := &Router{
			nil,
			StubDeployments{""},
			nil,
		}

		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		assertSuccessStatus(t, response)

		var got map[string]string
		if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}

		want := map[string]string{
			"id": "1ed4fdb9-2934-406f-a2bc-0e7cd8f301ae",
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})

	t.Run("DELETE / gets deployments error", func(t *testing.T) {
		mockLogger := NewMockLogger()
		controller := &Router{
			nil,
			StubDeployments{"Stub error"},
			&mockLogger,
		}

		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		assertServerError(t, response)
		mockLogger.assertErrors(t, []error{errors.New("Stub error")})
	})

	t.Run("GET /version returns version and sha1", func(t *testing.T) {
		controller := &Router{
			&StubBuild{
				"oogabooga",
				"boogaooga",
			},
			nil,
			nil,
		}

		request, _ := http.NewRequest(http.MethodGet, "/version", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		var got map[string]string

		if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		want := map[string]string{
			"sha1":    "oogabooga",
			"version": "boogaooga",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

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
