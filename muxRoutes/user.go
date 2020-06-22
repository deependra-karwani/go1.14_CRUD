package muxRoutes

import (
	"CRUD/controllers"
	"CRUD/workers"
	"net/http"

	"github.com/gorilla/mux"
)

func AddUserHandler(r *mux.Router) {
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.Register, w, r, done)
		<-done
	}).Methods(http.MethodPost)

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.Login, w, r, done)
		<-done
	}).Methods(http.MethodPut)

	r.HandleFunc("/forgot", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.ForgotPassword, w, r, done)
		<-done
	}).Methods(http.MethodPut)

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.Logout, w, r, done)
		<-done
	}).Methods(http.MethodGet)

	r.HandleFunc("/getAll", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.GetAllUsers, w, r, done)
		<-done
	}).Methods(http.MethodGet)

	r.HandleFunc("/getDetails", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.GetUserDetails, w, r, done)
		<-done
	}).Methods(http.MethodGet)

	r.HandleFunc("/updProf", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.UpdateUserProfile, w, r, done)
		<-done
	}).Methods(http.MethodPut)

	r.HandleFunc("/delAcc", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		defer close(done)
		workers.AddToQueue(controllers.DeleteUserAccount, w, r, done)
		<-done
	}).Methods(http.MethodDelete)
}
