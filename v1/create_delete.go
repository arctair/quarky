package v1

type CreateDelete interface {
	Create() (string, error)
	Delete() (string, error)
}
