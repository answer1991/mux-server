package route

import "net/http"

type Route interface {
	Method() (method string)
	Path() (path string)
	Process(w http.ResponseWriter, r *http.Request)
}
