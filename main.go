package main

import (
	"log"
	"net/http"
	storage "rpi-radio-alarm/helper"
	"rpi-radio-alarm/logging"
	"rpi-radio-alarm/resources/alarm"
	"rpi-radio-alarm/resources/radio"
	"rpi-radio-alarm/runner"
	"strconv"

	"github.com/gorilla/mux"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}

func main() {
	go runner.Runner()

	r := mux.NewRouter()
	r.HandleFunc("/health", health).Methods(http.MethodGet)

	alarm.SetUpRouter(r.PathPrefix("/alarm").Subrouter())
	radio.SetUpRouter(r.PathPrefix("/radio").Subrouter())

	storedData, err := storage.GetStoredData()
	if err != nil {
		logging.GetFatalLogger().Fatalln("error on getting stored data")
	}

	port := strconv.Itoa(storedData.Settings.Port)
	logging.GetInfoLogger().Printf("starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
