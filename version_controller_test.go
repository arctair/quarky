package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubBuild struct {
	sha1    string
	version string
}

func NewStubBuild(sha1, version string) *StubBuild {
	return &StubBuild{
		sha1,
		version,
	}
}

func (b *StubBuild) getSha1() string {
	return b.sha1
}

func (b *StubBuild) getVersion() string {
	return b.version
}

func TestVersionController(t *testing.T) {
	t.Run("GET returns version and sha1", func(t *testing.T) {
		versionController := NewVersionController(
			NewStubBuild(
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
