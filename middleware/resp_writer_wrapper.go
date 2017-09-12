package middleware

import "net/http"

type ResponseWriterWrapper struct {
	code int
	http.ResponseWriter
}

func (r *ResponseWriterWrapper) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *ResponseWriterWrapper) Write(content []byte) (int, error) {
	return r.ResponseWriter.Write(content)
}

func (r *ResponseWriterWrapper) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

func newResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		code:           200,
		ResponseWriter: w,
	}
}
