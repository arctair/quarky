package v1

import (
	"net/http"
)

// Router ...
type Router struct {
	rolloutsController *RolloutsController
	versionController  *VersionController
}

// NewRouter ...
func NewRouter(
	rolloutsController *RolloutsController,
	versionController *VersionController,
) *Router {
	return &Router{
		rolloutsController,
		versionController,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", router.rolloutsController.HandlerFunc())
	serveMux.Handle("/version", router.versionController.HandlerFunc())
	serveMux.ServeHTTP(w, request)
}
