package v1

import (
	"fmt"
)

type stubCreateDelete struct {
	err    error
	method string
}

// NewStubCreateDelete ...
func NewStubCreateDelete(method string, err error) CreateDelete {
	return &stubCreateDelete{err, method}
}

func (d stubCreateDelete) Create() error {
	if d.method != "Create" {
		return fmt.Errorf("Create called on %s stub", d.method)
	}
	return d.err
}

func (d stubCreateDelete) Delete() error {
	if d.method != "Delete" {
		return fmt.Errorf("Delete called on %s stub", d.method)
	}
	return d.err
}
