package v1

import (
	"net/http"
)

// RolloutsController ...
type RolloutsController struct {
	rollouts CreateDelete
	logger   Logger
}

// NewRolloutsController ...
func NewRolloutsController(r CreateDelete, l Logger) *RolloutsController {
	return &RolloutsController{r, l}
}

// HandlerFunc ...
func (c *RolloutsController) HandlerFunc() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				c.Post(w, r)
			case http.MethodDelete:
				c.Delete(w, r)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		},
	)
}

// Post ...
func (c *RolloutsController) Post(w http.ResponseWriter, r *http.Request) {
	if err := c.rollouts.Create(); err != nil {
		c.logger.error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Delete ...
func (c *RolloutsController) Delete(w http.ResponseWriter, r *http.Request) {
	if err := c.rollouts.Delete(); err != nil {
		c.logger.error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
