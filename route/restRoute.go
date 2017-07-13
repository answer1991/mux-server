package route

import (
	"encoding/json"
	"net/http"
)

type RestRoute interface {
	Method() string
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
			data, marshalErr := json.Marshal(map[string]interface{}{
				"err": err.Error.Error(),
			})

			if nil != marshalErr {
				w.WriteHeader(500)
				w.Write([]byte(marshalErr.Error()))
			} else {
				w.WriteHeader(err.Code)
				w.Write(data)
			}
		} else {
			data, marshalErr := json.Marshal(body)

			if nil != marshalErr {
				w.WriteHeader(500)
				w.Write([]byte(marshalErr.Error()))
			} else {
				w.Write(data)
			}
		}
	}
}
