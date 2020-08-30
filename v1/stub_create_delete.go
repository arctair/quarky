package v1

type stubCreateDelete struct {
	err error
}

// NewStubCreateDelete ...
func NewStubCreateDelete(stubError error) CreateDelete {
	return &stubCreateDelete{stubError}
}

func (d stubCreateDelete) Create() (string, error) {
	return "create", d.err
}

func (d stubCreateDelete) Delete() (string, error) {
	return "delete", d.err
}
