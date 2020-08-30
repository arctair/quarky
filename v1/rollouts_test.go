package v1

import (
	"errors"
	"testing"
)

func TestRollouts(t *testing.T) {
	t.Run("proxy Create to Deployments", func(t *testing.T) {
		rollouts := NewRollouts(NewStubCreateDelete(nil))

		id, err := rollouts.Create()
		if err != nil {
			t.Errorf("got %s want no error", err)
		}

		if id != "create" {
			t.Errorf("got %s want stub id 'create'", id)
		}
	})

	t.Run("proxy Create to Deployments returns error", func(t *testing.T) {
		rollouts := NewRollouts(NewStubCreateDelete(errors.New("stub error")))

		_, err := rollouts.Create()
		if err.Error() != "stub error" {
			t.Errorf("got %s want stub error", err)
		}
	})

	t.Run("proxy Delete to Deployments", func(t *testing.T) {
		rollouts := NewRollouts(NewStubCreateDelete(nil))

		id, err := rollouts.Delete()
		if err != nil {
			t.Errorf("got %s want no error", err)
		}

		if id != "delete" {
			t.Errorf("got %s want stub id 'delete'", id)
		}
	})

	t.Run("proxy Delete to Deployments returns error", func(t *testing.T) {
		rollouts := NewRollouts(NewStubCreateDelete(errors.New("stub error")))

		_, err := rollouts.Delete()
		if err.Error() != "stub error" {
			t.Errorf("got %s want stub error", err)
		}
	})
}
