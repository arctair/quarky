package v1

import (
	"encoding/json"
	"net/http"
)

// VersionController ...
type VersionController struct {
	build *Build
}

// NewVersionController ...
func NewVersionController(b *Build) *VersionController {
	return &VersionController{b}
}

// HandlerFunc ...
func (c *VersionController) HandlerFunc() http.Handler {
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
