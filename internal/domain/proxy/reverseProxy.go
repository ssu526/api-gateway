package proxy

import (
	"net/http"
)

type ReverseProxy interface {
	Forward(http.ResponseWriter, *http.Request, string)
}
