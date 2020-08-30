package v1

import (
	"errors"
	"testing"
)

func TestRollouts(t *testing.T) {
	t.Run("proxy Create to Deployments", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Create", nil),
		)

		if err := rollouts.Create(); err != nil {
			t.Errorf("got %s want no error", err)
		}
	})

	t.Run("proxy Create to Deployments returns error", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Create", errors.New("stub error")),
		)

		if err := rollouts.Create(); err.Error() != "stub error" {
			t.Errorf("got %s want stub error", err)
		}
	})

	t.Run("proxy Delete to Deployments", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Delete", nil),
		)

		if err := rollouts.Delete(); err != nil {
			t.Errorf("got %s want no error", err)
		}
	})

	t.Run("proxy Delete to Deployments returns error", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Delete", errors.New("stub error")),
		)

		if err := rollouts.Delete(); err == nil || err.Error() != "stub error" {
			t.Errorf("got %s want stub error", err)
		}
	})
}
