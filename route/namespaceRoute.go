package route

type NamespaceRoute struct {
	Namespace string

	Routes     []Route
	RestRoutes []RestRoute
}
