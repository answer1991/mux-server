package route

import "net/http"

type Route interface {
	Methods() []string
	Path() string
	Process(w http.ResponseWriter, r *http.Request)
}
