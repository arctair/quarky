package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"arctair.com/quarky/assertutil"
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
		controller := NewRolloutsController(
			NewStubCreateDelete("Create", nil),
			nil,
		)

		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		controller.HandlerFunc().ServeHTTP(response, request)

		assertutil.SuccessStatus(t, response.Result())
	})

	t.Run("POST when create rollout fails", func(t *testing.T) {
		mockLogger := NewMockLogger()
		controller := NewRolloutsController(
			NewStubCreateDelete("Create", errors.New("Stub error")),
			&mockLogger,
		)

		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		controller.HandlerFunc().ServeHTTP(response, request)

		assertutil.ServerErrorStatus(t, response.Result())
		mockLogger.assertErrors(t, []error{errors.New("Stub error")})
	})

	t.Run("DELETE deletes rollout", func(t *testing.T) {
		controller := NewRolloutsController(
			NewStubCreateDelete("Delete", nil),
			nil,
		)

		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		controller.HandlerFunc().ServeHTTP(response, request)

		assertutil.SuccessStatus(t, response.Result())
	})

	t.Run("DELETE when delete rollout fails", func(t *testing.T) {
		mockLogger := NewMockLogger()
		controller := NewRolloutsController(
			NewStubCreateDelete("Delete", errors.New("Stub error")),
			&mockLogger,
		)

		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		controller.HandlerFunc().ServeHTTP(response, request)

		assertutil.ServerErrorStatus(t, response.Result())
		mockLogger.assertErrors(t, []error{errors.New("Stub error")})
	})
}
