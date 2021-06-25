package main

import (
	"log"
	"net/http"
	"strconv"

	storage "github.com/bb4L/rpi-radio-alarm-go/helper"
	"github.com/bb4L/rpi-radio-alarm-go/logging"
	"github.com/bb4L/rpi-radio-alarm-go/resources/alarm"
	"github.com/bb4L/rpi-radio-alarm-go/resources/radio"
	"github.com/bb4L/rpi-radio-alarm-go/runner"

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
