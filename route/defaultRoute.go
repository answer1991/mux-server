package route

import "net/http"

type DefaultRoute interface {
	Process(w http.ResponseWriter, r *http.Request)
}
