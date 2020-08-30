package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubDeployments struct {
	stubError string
}

func NewStubDeployments(stubError string) *StubDeployments {
	return &StubDeployments{stubError}
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

func TestDeploymentsController(t *testing.T) {
	t.Run("POST creates deployment", func(t *testing.T) {
		deploymentsController := NewDeploymentsController(
			NewStubDeployments(""),
			nil,
		)

		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		deploymentsController.HandlerFunc().ServeHTTP(response, request)

		assertSuccessStatus(t, response)

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

	t.Run("POST gets deployments error", func(t *testing.T) {
		mockLogger := NewMockLogger()
		deploymentsController := NewDeploymentsController(
			NewStubDeployments("Stub error"),
			&mockLogger,
		)

		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		deploymentsController.HandlerFunc().ServeHTTP(response, request)

		assertServerError(t, response)
		mockLogger.assertErrors(t, []error{errors.New("Stub error")})
	})

	t.Run("DELETE deletes deployment", func(t *testing.T) {
		deploymentsController := NewDeploymentsController(
			NewStubDeployments(""),
			nil,
		)

		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		deploymentsController.HandlerFunc().ServeHTTP(response, request)

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

	t.Run("DELETE gets deployments error", func(t *testing.T) {
		mockLogger := NewMockLogger()
		deploymentsController := NewDeploymentsController(
			NewStubDeployments("Stub error"),
			&mockLogger,
		)

		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		deploymentsController.HandlerFunc().ServeHTTP(response, request)

		assertServerError(t, response)
		mockLogger.assertErrors(t, []error{errors.New("Stub error")})
	})
}
