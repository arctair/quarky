package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"arctair.com/quarky/testutil"
)

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

func TestRolloutsController(t *testing.T) {
	t.Run("POST creates rollout", func(t *testing.T) {
		rolloutController := NewRolloutController(
			NewStubCreateDelete(nil),
			nil,
		)

		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		rolloutController.HandlerFunc().ServeHTTP(response, request)

		testutil.AssertSuccessStatus(t, response)

		var got map[string]string
		if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}

		want := map[string]string{
			"id": "create",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("POST when create rollout fails", func(t *testing.T) {
		mockLogger := NewMockLogger()
		rolloutController := NewRolloutController(
			NewStubCreateDelete(errors.New("Stub error")),
			&mockLogger,
		)

		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		rolloutController.HandlerFunc().ServeHTTP(response, request)

		testutil.AssertServerError(t, response)
		mockLogger.assertErrors(t, []error{errors.New("Stub error")})
	})

	t.Run("DELETE deletes rollout", func(t *testing.T) {
		rollouterController := NewRolloutController(
			NewStubCreateDelete(nil),
			nil,
		)

		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		rollouterController.HandlerFunc().ServeHTTP(response, request)

		testutil.AssertSuccessStatus(t, response)

		var got map[string]string
		if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}

		want := map[string]string{
			"id": "delete",
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})

	t.Run("DELETE when delete rollout fails", func(t *testing.T) {
		mockLogger := NewMockLogger()
		rolloutController := NewRolloutController(
			NewStubCreateDelete(errors.New("Stub error")),
			&mockLogger,
		)

		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		rolloutController.HandlerFunc().ServeHTTP(response, request)

		testutil.AssertServerError(t, response)
		mockLogger.assertErrors(t, []error{errors.New("Stub error")})
	})
}
