package structs

import "net/http"

type API struct {
	Work func(w http.ResponseWriter, r *http.Request, done chan<- bool)
	W    http.ResponseWriter
	R    *http.Request
	Done chan<- bool
}
