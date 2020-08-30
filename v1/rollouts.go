package v1

// Rollouts ...
type Rollouts struct {
	deployments CreateDelete
}

// NewRollouts ...
func NewRollouts(d CreateDelete) CreateDelete {
	return &Rollouts{d}
}

// Create ...
func (r *Rollouts) Create() error {
	return r.deployments.Create()
}

// Delete ...
func (r *Rollouts) Delete() error {
	return r.deployments.Delete()
}
