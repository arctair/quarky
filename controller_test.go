package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestController(t *testing.T) {
	controller := &controller{}
	t.Run("GET /version returns scenario", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/version", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		var got map[string]string

		if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		want := map[string]string{"version": "failing acceptance tests"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
