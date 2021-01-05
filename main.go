package main

import (
	"log"
	"net/http"
	"rpi-radio-alarm/resources/alarm"
	"rpi-radio-alarm/resources/radio"

	"github.com/gorilla/mux"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	// case "POST":
	// 	w.WriteHeader(http.StatusCreated)
	// 	w.Write([]byte(`{"message": "post called"}`))
	// case "PUT":
	// 	w.WriteHeader(http.StatusAccepted)
	// 	w.Write([]byte(`{"message": "put called"}`))
	// case "DELETE":
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`{"message": "delete called"}`))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	go alarm.Runner()

	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	alarm.SetUpRouter(apiV1.PathPrefix("/alarm").Subrouter())
	radio.SetUpRouter(apiV1.PathPrefix("/radio").Subrouter())

	r.HandleFunc("/health", health).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}
