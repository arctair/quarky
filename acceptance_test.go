// +build acceptance

package main_test

import (
	"bufio"
	"encoding/json"
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

	"arctair.com/quarky/assertutil"
	"github.com/cenkalti/backoff/v4"
)

func TestAcceptance(t *testing.T) {
	clusterUrl := os.Getenv("CLUSTER_URL")
	if clusterUrl == "" {
		fmt.Println("Please set environment variable CLUSTER_URL. It's the URL of some ingress that can route to an instance of Quarky. Try the output of scripts/getLoadBalancerUrl.")
		os.Exit(1)
	}
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:5000"

		assertutil.NotError(t, exec.Command("sh", "build").Run())

		command := exec.Command("bin/quarky")
		stderr, err := command.StderrPipe()
		assertutil.NotError(t, err)
		assertutil.NotError(t, command.Start())
		defer dumpPipe("app:", stderr)
		defer command.Process.Kill()

		assertutil.NotError(
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

	t.Run("new change is made active", func(t *testing.T) {
		client := &http.Client{}

		t.Cleanup(func() {
			request, err := http.NewRequest("DELETE", baseUrl, nil)
			assertutil.NotError(t, err)
			response, err := client.Do(request)
			assertutil.NotError(t, err)
			assertutil.SuccessStatus(t, response)

			request, err = http.NewRequest("GET", clusterUrl, nil)
			assertutil.NotError(t, err)
			request.Header.Add("Host", "quarky-test")

			assertutil.NotError(
				t,
				backoff.Retry(
					func() error {
						response, err = client.Do(request)
						if err != nil {
							return err
						}
						if response.StatusCode < 500 {
							return fmt.Errorf(
								"got %s want server error indicating deployment was torn down",
								response.Status,
							)
						}
						return nil
					},
					NewExponentialBackOff(5*time.Second),
				),
			)
		})

		response, err := http.Post(baseUrl, "", nil)
		assertutil.NotError(t, err)
		assertutil.SuccessStatus(t, response)

		request, err := http.NewRequest("GET", clusterUrl, nil)
		assertutil.NotError(t, err)
		request.Header.Add("Host", "quarky-test")
		assertutil.NotError(
			t,
			backoff.Retry(
				func() error {
					response, err = client.Do(request)
					if err != nil {
						return err
					}
					if response.StatusCode > 299 {
						return fmt.Errorf(
							"got %s want successful status",
							response.Status,
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
		assertutil.NotError(t, err)

		want := map[string]string{
			"scenario": "passing acceptance tests",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("GET /version returns sha1 and version", func(t *testing.T) {
		response, err := http.Get(fmt.Sprintf("%s/version", baseUrl))
		assertutil.NotError(t, err)

		var got map[string]string
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&got)
		assertutil.NotError(t, err)

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
