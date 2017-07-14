package server

import (
	"context"
	"fmt"
	"github.com/answer1991/mux-server/route"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

func NewServer(port int) *Server {
	return &Server{
		Port:            port,
		routes:          []route.Route{},
		restRoutes:      []route.RestRoute{},
		namespaceRoutes: []*route.NamespaceRoute{},
	}
}

type Server struct {
	Port    int
	Version string

	routes          []route.Route
	restRoutes      []route.RestRoute
	namespaceRoutes []*route.NamespaceRoute

	muxRouter *mux.Router
	staticDir string
}

func (this *Server) Serve(ctx context.Context) (err error) {
	this.init()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", this.Port))

	if nil != err {
		return err
	}

	go func() {
		http.Serve(l, this.muxRouter)
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				l.Close()
			}
		}
	}()

	return nil
}

func (this *Server) init() {
	w := mux.NewRouter()

	if "" == this.Version {
		this.Version = "v{version:[0-9.]+}"
	}

	for _, r := range this.routes {
		w.
			PathPrefix(fmt.Sprintf("/%s", this.Version)).
			Path(r.Path()).
			Methods(r.Method()).
			HandlerFunc(r.Process)

		w.
			Path(r.Path()).
			Methods(r.Method()).
			HandlerFunc(r.Process)
	}

	for _, r := range this.restRoutes {
		w.
			PathPrefix(fmt.Sprintf("/%s", this.Version)).
			Path(r.Path()).
			Methods(r.Method()).
			HandlerFunc(route.ConvertToHandlerFunc(r.Process))

		w.
			Path(r.Path()).
			Methods(r.Method()).
			HandlerFunc(route.ConvertToHandlerFunc(r.Process))
	}

	for _, nr := range this.namespaceRoutes {
		for _, r := range nr.Routes {
			w.
				PathPrefix(fmt.Sprintf("/%s/%s", this.Version, nr.Namespace)).
				Path(r.Path()).
				Methods(r.Method()).
				HandlerFunc(r.Process)

			w.
				PathPrefix(fmt.Sprintf("/%s", nr.Namespace)).
				Path(r.Path()).
				Methods(r.Method()).
				HandlerFunc(r.Process)
		}

		for _, r := range nr.RestRoutes {
			w.
				PathPrefix(fmt.Sprintf("/%s/%s", this.Version, nr.Namespace)).
				Path(r.Path()).
				Methods(r.Method()).
				HandlerFunc(route.ConvertToHandlerFunc(r.Process))

			w.
				PathPrefix(fmt.Sprintf("/%s", nr.Namespace)).
				Path(r.Path()).
				Methods(r.Method()).
				HandlerFunc(route.ConvertToHandlerFunc(r.Process))
		}
	}

	if "" != this.staticDir {
		w.
			PathPrefix("/").
			Handler(http.StripPrefix("/", http.FileServer(http.Dir(this.staticDir))))
	}

	this.muxRouter = w
}

func (this *Server) AddRoute(r route.Route) {
	this.routes = append(this.routes, r)
}

func (this *Server) AddRestRoute(r route.RestRoute) {
	this.restRoutes = append(this.restRoutes, r)
}

func (this *Server) AddNamespaceRoute(r *route.NamespaceRoute) {
	this.namespaceRoutes = append(this.namespaceRoutes, r)
}

func (this *Server) SetStaticFilePath(dir string) {
	this.staticDir = dir
}
