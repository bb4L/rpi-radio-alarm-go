package radio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	storage "rpi-radio-alarm/helper"

	"github.com/gorilla/mux"
)

// GetRadio returns the radio
func getRadio(w http.ResponseWriter, r *http.Request) {
	var err error

	storedData, err := storage.GetStoredData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on getting stored data"}`))
		return
	}
	jsonData, err := json.Marshal(storedData.Radio)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on creating json data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	return
}

type changeValue struct {
	ChangeValue string `json:"switch"`
}

func postRadio(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var msg changeValue
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	storedData, err := storage.GetStoredData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on getting stored data"}`))
		return
	}

	if !storedData.Radio.Running && msg.ChangeValue == "on" {
		storedData.Radio.StartRadio()
	}

	if storedData.Radio.Running && msg.ChangeValue == "off" {
		storedData.Radio.StopRadio()
	}

	jsonData, err := json.Marshal(storedData.Radio)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on creating json data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	storage.SaveData(storedData)
}

// SetUpRouter set up router for radio endpoints
func SetUpRouter(r *mux.Router) {
	r.HandleFunc("", postRadio).Methods(http.MethodPost)
	r.HandleFunc("", getRadio).Methods(http.MethodGet)
}
