package radio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bb4L/rpi-radio-alarm-go-library/logging"
	"github.com/bb4L/rpi-radio-alarm-go-library/storage"
	"github.com/bb4L/rpi-radio-alarm-go/constants"

	"github.com/gorilla/mux"
)

var logger = logging.GetLogger(os.Stdout, constants.DefaultPrefix, "apiserver", "radio")

var storageHelper *storage.Helper

// SetUpRouter set up router for radio endpoints
func SetUpRouter(r *mux.Router, storageH *storage.Helper) {
	logger.Println("setup router")

	storageHelper = storageH

	r.HandleFunc("", postRadio).Methods(http.MethodPost)
	r.HandleFunc("", getRadio).Methods(http.MethodGet)
}

// GetRadio returns the radio
func getRadio(w http.ResponseWriter, r *http.Request) {
	radio, err := storageHelper.GetRadio(false)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on getting stored data"}`))
		return
	}
	jsonData, err := json.Marshal(radio)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on creating json data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

type changeValue struct {
	ChangeValue string `json:"switch"`
}

func postRadio(w http.ResponseWriter, r *http.Request) {
	radio, err := storageHelper.GetRadio(false)
	defer storageHelper.SaveRadio(radio)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on getting stored data"}`))
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on reading body"}`))
		return
	}
	defer r.Body.Close()

	var msg changeValue
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if !radio.Running && msg.ChangeValue == "on" {
		radio.StartRadio()
	}

	if radio.Running && msg.ChangeValue == "off" {
		radio.StopRadio()
	}

	jsonData, err := json.Marshal(radio)
	if err != nil {
		logger.Printf("json error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on creating json data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}
