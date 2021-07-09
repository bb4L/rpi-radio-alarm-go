package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bb4L/rpi-radio-alarm-go-library/logging"
	"github.com/bb4L/rpi-radio-alarm-go-library/storage"
	"github.com/bb4L/rpi-radio-alarm-go/apiserver/alarm"
	"github.com/bb4L/rpi-radio-alarm-go/apiserver/radio"
	"github.com/bb4L/rpi-radio-alarm-go/constants"
	"github.com/gorilla/mux"
)

var logger = logging.GetLogger(os.Stdout, constants.DefaultPrefix, "apiserver")

// StartAPIServer starts the http server to interact with the data
func StartAPIServer(storageHelper *storage.Helper) {

	logger.Println("starting api server")
	r := mux.NewRouter()
	r.HandleFunc("/health", health).Methods(http.MethodGet)

	alarm.SetUpRouter(r.PathPrefix("/alarm").Subrouter(), storageHelper)
	radio.SetUpRouter(r.PathPrefix("/radio").Subrouter(), storageHelper)

	settings, err := storageHelper.GetSettings()
	if err != nil {
		logger.Fatalln("error on getting stored data")
	}
	port := strconv.Itoa(settings.Port)

	logger.Printf("starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}
