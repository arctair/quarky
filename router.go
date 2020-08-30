package main

import (
	"net/http"
)

// Router ...
type Router struct {
	deploymentsController *DeploymentsController
	versionController     *VersionController
}

// NewRouter ...
func NewRouter(
	deploymentsController *DeploymentsController,
	versionController *VersionController,
) Router {
	return Router{
		deploymentsController,
		versionController,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", router.deploymentsController.HandlerFunc())
	serveMux.Handle("/version", router.versionController.HandlerFunc())
	serveMux.ServeHTTP(w, request)
}
