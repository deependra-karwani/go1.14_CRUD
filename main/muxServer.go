package main

import (
	"CRUD/middlewares"
	"CRUD/muxRoutes"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	userSr := r.PathPrefix("/").Subrouter()
	fileSr := r.PathPrefix("/").Subrouter()

	userSr.Use(middlewares.UserAuthMux)

	muxRoutes.AddUserHandler(userSr)
	muxRoutes.AddFileHandler(fileSr)

	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Accept", "Origin", "token"})
	headersOK := handlers.ExposedHeaders([]string{"token"})

	CORSServer := handlers.CORS(originsOk, methodsOk, headersOk, headersOK)(r)
	if err := http.ListenAndServe(":8080", CORSServer); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
