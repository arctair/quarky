package main_test

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
)

func TestAcceptance(t *testing.T) {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:5000"

		assertNotError(t, exec.Command("sh", "build").Run())

		command := exec.Command("bin/quarky")
		stderr, err := command.StderrPipe()
		assertNotError(t, err)
		assertNotError(t, command.Start())
		defer dumpPipe("app:", stderr)
		defer command.Process.Kill()

		assertNotError(
			t,
			backoff.Retry(
				func() error {
					_, err := http.Get(baseUrl)
					return err
				},
				NewExponentialBackOff(3*time.Second),
			),
		)
	}

	t.Run("POST / deploys quarky-test", func(t *testing.T) {
		client := &http.Client{}

		t.Cleanup(func() {
			request, err := http.NewRequest("DELETE", baseUrl, nil)
			assertNotError(t, err)
			response, err := client.Do(request)
			assertNotError(t, err)

			if response.StatusCode > 299 {
				t.Fatalf("got %s want successful status", response.Status)
			}

			request, err = http.NewRequest("GET", "http://172.17.0.3", nil)
			assertNotError(t, err)
			request.Header.Add("Host", "quarky-test")

			assertNotError(
				t,
				backoff.Retry(
					func() error {
						response, err = client.Do(request)
						if err != nil {
							return err
						}
						if response.StatusCode < 500 {
							return errors.New(
								fmt.Sprintf(
									"got %s want server error indicating deployment was torn down",
									response.Status,
								),
							)
						}
						return nil
					},
					NewExponentialBackOff(5*time.Second),
				),
			)
		})
		response, err := http.Post(baseUrl, "", nil)
		assertNotError(t, err)

		if response.StatusCode > 299 {
			t.Fatalf("got %s want successful status", response.Status)
		}

		request, err := http.NewRequest("GET", "http://172.17.0.3", nil)
		assertNotError(t, err)
		request.Header.Add("Host", "quarky-test")
		assertNotError(
			t,
			backoff.Retry(
				func() error {
					response, err = client.Do(request)
					if err != nil {
						return err
					}
					if response.StatusCode > 299 {
						return errors.New(
							fmt.Sprintf(
								"got %s want successful status",
								response.Status,
							),
						)
					}
					return nil
				},
				NewExponentialBackOff(5*time.Second),
			),
		)

		var got map[string]string
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&got)
		assertNotError(t, err)

		want := map[string]string{
			"scenario": "passing acceptance tests",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("GET /version returns sha1 and version", func(t *testing.T) {
		response, err := http.Get(fmt.Sprintf("%s/version", baseUrl))
		assertNotError(t, err)

		var got map[string]string
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&got)
		assertNotError(t, err)

		sha1Pattern := regexp.MustCompile("^[0-9a-f]{40}(-dirty)?$")
		versionPattern := regexp.MustCompile("^\\d+\\.\\d+\\.\\d+$")

		if !sha1Pattern.MatchString(got["sha1"]) {
			t.Errorf("got sha1 %s want 40 hex digits", got["sha1"])
		}
		if !versionPattern.MatchString(got["version"]) && !sha1Pattern.MatchString(got["version"]) {
			t.Errorf("got version %s want semver or 40 hex digits", got["version"])
		}
	})

}

func assertNotError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func dumpPipe(prefix string, p io.ReadCloser) {
	s := bufio.NewScanner(p)
	for s.Scan() {
		log.Printf("%s: %s", prefix, s.Text())
	}
	if err := s.Err(); err != nil {
		log.Printf("Failed to dump pipe: %s", err)
	}
}

func NewExponentialBackOff(timeout time.Duration) *backoff.ExponentialBackOff {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      timeout,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	return b
}
