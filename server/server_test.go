package server

import (
	"context"
	"github.com/answer1991/mux-server/route"
	"log"
	"net/http"
	"path"
	"testing"
)

type testStruct struct {
	Name string `json:"name"`
}

type testRouter struct {
	Context string
}

func (r *testRouter) Methods() []string {
	return []string{http.MethodGet}
}

func (r *testRouter) Path() string {
	return "/test"
}

func (r *testRouter) Process(req *http.Request) (body interface{}, error *route.HttpServerError) {
	test := &testStruct{}
	UnmarshalRequestBodyToJson(req, test)

	return map[string]interface{}{
		"test":    "hello-world",
		"value":   req.FormValue("name"),
		"body":    test,
		"context": r.Context,
	}, nil
}

//var TestRestRoute = &route.RestRoute{
//	Method: http.MethodPost,
//	Path:   "/test",
//	Process: func(r *http.Request) (body interface{}, error *route.HttpServerError) {
//		test := &testStruct{}
//		UnmarshalRequestBodyToJson(r, test)
//
//		return map[string]interface{}{
//			"test":  "hello-world",
//			"value": r.FormValue("name"),
//			"body":  test,
//		}, nil
//	},
//}

//type TestRestErrRoute struct {
//}
//
//func (this *TestRestErrRoute) Method() (ret string) {
//	return http.MethodGet
//}
//
//func (this *TestRestErrRoute) Path() (ret string) {
//	return "/testErr"
//}
//
//func (this *TestRestErrRoute) Process(r *http.Request) (body interface{}, error *route.HttpServerError) {
//	return nil, &route.HttpServerError{
//		Code:  503,
//		Error: errors.New("something goes wrong"),
//	}
//}

func TestNewServer(t *testing.T) {
	s := NewServer(80)

	//s.Version = "v1"

	s.AddRestRoute(&testRouter{
		Context: "hello world",
	})

	s.SetStaticFilePath(path.Join("..", "public"))

	cxt := context.Background()

	log.Println(s.Serve(cxt))

	<-make(chan string)
}
