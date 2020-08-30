package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestVersionController(t *testing.T) {
	t.Run("GET returns version and sha1", func(t *testing.T) {
		versionController := NewVersionController(
			NewBuild(
				"oogabooga",
				"boogaooga",
			),
		)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		versionController.HandlerFunc().ServeHTTP(response, request)

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
