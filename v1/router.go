package v1

import (
	"net/http"
)

// Router ...
type Router struct {
	deploymentsController *RolloutController
	versionController     *VersionController
}

// NewRouter ...
func NewRouter(
	deploymentsController *RolloutController,
	versionController *VersionController,
) *Router {
	return &Router{
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
