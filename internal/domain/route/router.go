package route

import "net/http"

type Router interface {
	FindRoute(*http.Request) (*Route, bool)
}
