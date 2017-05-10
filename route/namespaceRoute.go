package route

type NamespaceRoute interface {
	Namespace() (method string)

	Routes() []Route
	RestRoutes() []RestRoute
}
