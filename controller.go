package main

import (
	"encoding/json"
	"net/http"
)

type controller struct {
	build       Build
	deployments Deployments
	logger      Logger
}

func (c *controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()
	deploymentsController := DeploymentsController{
		c.deployments,
		c.logger,
	}

	router.Handle(
		"/version",
		getVersion(c),
	)

	router.Handle(
		"/",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					deploymentsController.Post(w, r)
				case http.MethodDelete:
					deploymentsController.Delete(w, r)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			},
		),
	)

	router.ServeHTTP(w, r)
}

// DeploymentsController ...
type DeploymentsController struct {
	Deployments Deployments
	Logger      Logger
}

// Post ...
func (c *DeploymentsController) Post(w http.ResponseWriter, r *http.Request) {
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
func (c *DeploymentsController) Delete(w http.ResponseWriter, r *http.Request) {
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

func getVersion(c *controller) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			bytes, err := json.Marshal(
				map[string]string{
					"sha1":    c.build.getSha1(),
					"version": c.build.getVersion(),
				},
			)
			if err != nil {
				panic(err)
			}
			w.Write(bytes)
		},
	)
}
