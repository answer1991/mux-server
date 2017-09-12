package route

import (
	"encoding/json"
	"net/http"
)

type RestRoute interface {
	Methods() []string
	Path() string
	Process(r *http.Request) (body interface{}, error *HttpServerError)
}

const (
	contentTypeKey  = "Content-Type"
	contentTypeJson = "application/json"
)

func ConvertToHandlerFunc(process func(r *http.Request) (body interface{}, error *HttpServerError)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := process(r)

		w.Header().Add(contentTypeKey, contentTypeJson)

		if nil != err {
			http.Error(w, err.Error.Error(), err.Code)
		} else {
			if data, ok := body.([]byte); ok {
				w.Write(data)
			} else {
				data, marshalErr := json.Marshal(body)

				if nil != marshalErr {
					http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
				} else {
					w.Write(data)
				}
			}
		}
	}
}
