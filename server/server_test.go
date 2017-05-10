package server

import (
	"errors"
	"github.com/answer1991/mux-server/route"
	"log"
	"net/http"
	"testing"
)

type TestRestRoute struct {
}

func (this *TestRestRoute) Method() (ret string) {
	return http.MethodPost
}

func (this *TestRestRoute) Path() (ret string) {
	return "/test"
}

func (this *TestRestRoute) Process(r *http.Request) (body interface{}, error *route.HttpServerError) {
	return map[string]interface{}{
		"test":  "hello-world",
		"value": r.FormValue("name"),
		"body":  r.PostFormValue(""),
	}, nil
}

type TestRestErrRoute struct {
}

func (this *TestRestErrRoute) Method() (ret string) {
	return http.MethodGet
}

func (this *TestRestErrRoute) Path() (ret string) {
	return "/testErr"
}

func (this *TestRestErrRoute) Process(r *http.Request) (body interface{}, error *route.HttpServerError) {
	return nil, &route.HttpServerError{
		Code:  503,
		Error: errors.New("something goes wrong"),
	}
}

func TestNewServer(t *testing.T) {
	s := NewServer(80)

	s.Version = "v1"

	s.AddRestRoute(&TestRestRoute{})
	s.AddRestRoute(&TestRestErrRoute{})

	log.Fatal(s.Serve())
}
