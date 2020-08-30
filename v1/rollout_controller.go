package v1

import (
	"encoding/json"
	"net/http"
)

// RolloutController ...
type RolloutController struct {
	Deployments Deployments
	Logger      Logger
}

// NewRolloutController ...
func NewRolloutController(d Deployments, l Logger) *RolloutController {
	return &RolloutController{d, l}
}

// HandlerFunc ...
func (c *RolloutController) HandlerFunc() http.Handler {
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
func (c *RolloutController) Post(w http.ResponseWriter, r *http.Request) {
	id, err := c.Deployments.Create()
	if err != nil {
		c.Logger.error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(
		map[string]string{
			"id": id,
		},
	)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
}

// Delete ...
func (c *RolloutController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := c.Deployments.Delete()
	if err != nil {
		c.Logger.error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(
		map[string]string{
			"id": id,
		},
	)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
}
