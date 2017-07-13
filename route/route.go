package route

import "net/http"

type Route interface {
	Method() string
	Path() string
	Process(w http.ResponseWriter, r *http.Request)
}
