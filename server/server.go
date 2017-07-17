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

	defaultRoute route.DefaultRoute

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

func (this *Server) initRoute(r route.Route, namespace string) {
	path := r.Path()

	if "" != namespace {
		path = fmt.Sprintf("/%s%s", namespace, path)
	}

	vRouter := this.muxRouter.
		PathPrefix(fmt.Sprintf("/%s", this.Version)).
		Path(r.Path())

	router := this.muxRouter.
		Path(r.Path())

	if methods := r.Methods(); nil != methods && 0 < len(methods) {
		vRouter.Methods(methods...)
		router.Methods(methods...)
	}

	vRouter.HandlerFunc(r.Process)
	router.HandlerFunc(r.Process)
}

func (this *Server) initRestRoute(r route.RestRoute, namespace string) {
	path := r.Path()

	if "" != namespace {
		path = fmt.Sprintf("/%s%s", namespace, path)
	}

	vRouter := this.muxRouter.
		PathPrefix(fmt.Sprintf("/%s", this.Version)).
		Path(r.Path())

	router := this.muxRouter.
		Path(r.Path())

	if methods := r.Methods(); nil != methods && 0 < len(methods) {
		vRouter.Methods(methods...)
		router.Methods(methods...)
	}

	fn := route.ConvertToHandlerFunc(r.Process)

	vRouter.HandlerFunc(fn)
	router.HandlerFunc(fn)
}

func (this *Server) init() {
	w := mux.NewRouter()
	this.muxRouter = w

	if "" == this.Version {
		this.Version = "v{version:[0-9.]+}"
	}

	for _, r := range this.routes {
		this.initRoute(r, "")
	}

	for _, r := range this.restRoutes {
		this.initRestRoute(r, "")
	}

	for _, nr := range this.namespaceRoutes {
		for _, r := range nr.Routes {
			this.initRoute(r, nr.Namespace)
		}

		for _, r := range nr.RestRoutes {
			this.initRestRoute(r, nr.Namespace)
		}
	}

	if "" != this.staticDir {
		w.
			PathPrefix("/").
			Handler(http.StripPrefix("/", http.FileServer(http.Dir(this.staticDir))))
	}

	if nil != this.defaultRoute {
		w.
			PathPrefix("/").
			HandlerFunc(this.defaultRoute.Process)
	}

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

func (this *Server) SetDefaultRoute(r route.DefaultRoute) {
	this.defaultRoute = r
}
