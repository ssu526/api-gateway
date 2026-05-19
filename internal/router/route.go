package router

type Route struct {
	Path        string
	ServiceName string
	Methods     []string
	Middlewares []string
}
