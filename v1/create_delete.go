package v1

// CreateDelete ...
type CreateDelete interface {
	Create() error
	Delete() error
}
