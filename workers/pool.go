package workers

import (
	"CRUD/structs"
	"net/http"
)

var (
	APIQueue   chan structs.API
	maxJobs    int
	maxWorkers int
)

func init() {
	maxWorkers = 11
	maxJobs = 1000
	APIQueue = make(chan structs.API, maxJobs)
	NewWorkers()
}

func Worker() {
	for {
		select {
		case api := <-APIQueue:
			go api.Work(api.W, api.R, api.Done)
		}
	}
}

func NewWorkers() {
	for i := 0; i < maxWorkers; i++ {
		go Worker()
	}
}

func AddToQueue(fnxn func(w http.ResponseWriter, r *http.Request, done chan<- bool), w http.ResponseWriter, r *http.Request, done chan<- bool) {
	APIQueue <- structs.API{fnxn, w, r, done}
}
