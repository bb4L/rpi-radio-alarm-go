package main

import (
	"log"
	"net/http"
	"rpi-radio-alarm/resources/alarm"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "get called"}`))
	case "POST":
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "post called"}`))
	case "PUT":
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "put called"}`))
	case "DELETE":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "delete called"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

func main() {

	// TODO: set up server
	// TODO: go routine to run the alarm check (while true.. sleep)

	go alarm.Runner()
	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	alarm.SetUpRouter(apiV1.PathPrefix("alarm").Subrouter())
	radio.SetUpRouter(apiV1.PathPrefix("radio").Subrouter())

	r.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", r))
}
