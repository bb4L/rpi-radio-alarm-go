package alarm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	storage "github.com/bb4L/rpi-radio-alarm-go/helper"

	alarmtypes "github.com/bb4L/rpi-radio-alarm-go-library/types"

	"strconv"

	"github.com/gorilla/mux"
)

// Set up router for alarm endpoints
func SetUpRouter(router *mux.Router) {
	router.HandleFunc("", getAlarms).Methods(http.MethodGet)
	router.HandleFunc("", addAlarm).Methods(http.MethodPost)

	router.HandleFunc("/{idx}", getAlarm).Methods(http.MethodGet)
	router.HandleFunc("/{idx}", changeAlarm).Methods(http.MethodPut)
	router.HandleFunc("/{idx}", deleteAlarm).Methods(http.MethodDelete)
}

func getAlarms(w http.ResponseWriter, r *http.Request) {
	storedData, err := storage.GetStoredData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on getting stored data"}`))
		return
	}

	jsonData, err := json.Marshal(storedData.Alarms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on creating json data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func addAlarm(w http.ResponseWriter, r *http.Request) {
	var err error

	storedData, err := storage.GetStoredData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on getting stored data"}`))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on reading body"}`))
		return
	}

	var alarm alarmtypes.Alarm

	err = json.Unmarshal(body, &alarm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on unmarshal body"}`))
		return
	}

	storedData.Alarms = append(storedData.Alarms, alarm)
	storage.SaveData(storedData)

	jsonData, err := json.Marshal(storedData.Alarms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on creating json data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getAlarm(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)

	alarmIdx := -1
	var err error

	if val, ok := pathParams["idx"]; ok {
		alarmIdx, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number as path param"}`))
			return
		}
	}

	storedData, err := storage.GetStoredData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on getting stored data"}`))
		return
	}

	jsonData, err := json.Marshal(storedData.Alarms[alarmIdx])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on creating json data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func changeAlarm(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)

	alarmIdx := -1
	var err error

	if val, ok := pathParams["idx"]; ok {
		alarmIdx, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number as path param"}`))
			return
		}
	}

	storedData, err := storage.GetStoredData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on getting stored data"}`))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on reading body"}`))
		return
	}

	var alarm alarmtypes.Alarm

	err = json.Unmarshal(body, &alarm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on unmarshal body"}`))
		return
	}

	storedData.Alarms[alarmIdx] = alarm
	storage.SaveData(storedData)

	jsonData, err := json.Marshal(alarm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on creating json data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteAlarm(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)

	alarmIdx := -1
	var err error

	if val, ok := pathParams["idx"]; ok {
		alarmIdx, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number as path param"}`))
			return
		}
	}

	storedData, err := storage.GetStoredData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on getting stored data"}`))
		return
	}

	storedData.Alarms = append(storedData.Alarms[:alarmIdx], storedData.Alarms[alarmIdx+1:]...)
	storage.SaveData(storedData)

	jsonData, err := json.Marshal(storedData.Alarms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error on creating json data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
