package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/answer1991/mux-server/middleware"
	"github.com/answer1991/mux-server/route"
	"github.com/gorilla/mux"
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

func (s *Server) Serve(ctx context.Context) (err error) {
	s.init()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))

	if nil != err {
		return err
	}

	go func() {
		http.Serve(l, s.muxRouter)
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

func (s *Server) initRoute(r route.Route, namespace string) {
	path := r.Path()

	if "" != namespace {
		path = fmt.Sprintf("/%s%s", namespace, path)
	}

	vRouter := s.muxRouter.
		PathPrefix(fmt.Sprintf("/%s", s.Version)).
		Path(path)

	router := s.muxRouter.
		Path(path)

	if methods := r.Methods(); nil != methods && 0 < len(methods) {
		vRouter.Methods(methods...)
		router.Methods(methods...)
	}

	vRouter.Handler(
		middleware.Use(http.HandlerFunc(r.Process), middleware.Access))
	router.Handler(
		middleware.Use(http.HandlerFunc(r.Process), middleware.Access))
}

func (s *Server) initRestRoute(r route.RestRoute, namespace string) {
	path := r.Path()

	if "" != namespace {
		path = fmt.Sprintf("/%s%s", namespace, path)
	}

	vRouter := s.muxRouter.
		PathPrefix(fmt.Sprintf("/%s", s.Version)).
		Path(path)

	router := s.muxRouter.
		Path(path)

	if methods := r.Methods(); nil != methods && 0 < len(methods) {
		vRouter.Methods(methods...)
		router.Methods(methods...)
	}

	fn := route.ConvertToHandlerFunc(r.Process)

	vRouter.Handler(
		middleware.Use(http.HandlerFunc(fn), middleware.Access))
	router.Handler(
		middleware.Use(http.HandlerFunc(fn), middleware.Access))
}

func (s *Server) init() {
	w := mux.NewRouter()
	s.muxRouter = w

	if "" == s.Version {
		s.Version = "v{version:[0-9.]+}"
	}

	for _, r := range s.routes {
		s.initRoute(r, "")
	}

	for _, r := range s.restRoutes {
		s.initRestRoute(r, "")
	}

	for _, nr := range s.namespaceRoutes {
		for _, r := range nr.Routes {
			s.initRoute(r, nr.Namespace)
		}

		for _, r := range nr.RestRoutes {
			s.initRestRoute(r, nr.Namespace)
		}
	}

	if "" != s.staticDir {
		w.
			PathPrefix("/").
			Handler(http.StripPrefix("/", http.FileServer(http.Dir(s.staticDir))))
	}

	if nil != s.defaultRoute {
		w.
			PathPrefix("/").
			HandlerFunc(s.defaultRoute.Process)
	} else {
		//w.
		//	PathPrefix("/").
		//	Handler(middleware.Use(http.NotFoundHandler(), middleware.Access))
	}

}

func (s *Server) AddRoute(r route.Route) {
	s.routes = append(s.routes, r)
}

func (s *Server) AddRestRoute(r route.RestRoute) {
	s.restRoutes = append(s.restRoutes, r)
}

func (s *Server) AddNamespaceRoute(r *route.NamespaceRoute) {
	s.namespaceRoutes = append(s.namespaceRoutes, r)
}

func (s *Server) SetStaticFilePath(dir string) {
	s.staticDir = dir
}

func (s *Server) SetDefaultRoute(r route.DefaultRoute) {
	s.defaultRoute = r
}
