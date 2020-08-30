package v1

// Rollouts ...
type Rollouts struct {
	deployments CreateDelete
	services    CreateDelete
}

// NewRollouts ...
func NewRollouts(
	d CreateDelete,
	s CreateDelete,
) CreateDelete {
	return &Rollouts{d, s}
}

// Create ...
func (r *Rollouts) Create() error {
	if err := r.deployments.Create(); err != nil {
		return err
	}
	return r.services.Create()
}

// Delete ...
func (r *Rollouts) Delete() error {
	if err := r.deployments.Delete(); err != nil {
		return err
	}
	return r.services.Delete()
}
