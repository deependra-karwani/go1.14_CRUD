package muxRoutes

import (
	"CRUD/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func AddFileHandler(r *mux.Router) {
	fileServer := http.FileServer(http.Dir("../images"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", middlewares.Neuter(fileServer)))
}
