package v1

import (
	"errors"
	"testing"
)

func TestRollouts(t *testing.T) {
	t.Run("create rollout creates deployment and service", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Create", nil),
			NewStubCreateDelete("Create", nil),
		)

		if err := rollouts.Create(); err != nil {
			t.Errorf("got %s want no error", err)
		}
	})

	t.Run("create rollout when create deployment has error", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Create", errors.New("stub error")),
			NewStubCreateDelete("Create", nil),
		)

		if err := rollouts.Create(); err == nil || err.Error() != "stub error" {
			t.Errorf("got %s want stub error", err)
		}
	})

	t.Run("create rollout when create service has error", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Create", nil),
			NewStubCreateDelete("Create", errors.New("stub error")),
		)

		if err := rollouts.Create(); err == nil || err.Error() != "stub error" {
			t.Errorf("got %s want stub error", err)
		}
	})

	t.Run("delete rollout deletes deployment and service", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Delete", nil),
			NewStubCreateDelete("Delete", nil),
		)

		if err := rollouts.Delete(); err != nil {
			t.Errorf("got %s want no error", err)
		}
	})

	t.Run("delete rollout when delete deployment has error", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Delete", errors.New("stub error")),
			NewStubCreateDelete("Delete", nil),
		)

		if err := rollouts.Delete(); err == nil || err.Error() != "stub error" {
			t.Errorf("got %s want stub error", err)
		}
	})

	t.Run("delete rollout when delete service has error", func(t *testing.T) {
		rollouts := NewRollouts(
			NewStubCreateDelete("Delete", nil),
			NewStubCreateDelete("Delete", errors.New("stub error")),
		)

		if err := rollouts.Delete(); err == nil || err.Error() != "stub error" {
			t.Errorf("got %s want stub error", err)
		}
	})
}
