package alarm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	storage "rpi-radio-alarm/helper"
	alarmtypes "rpi-radio-alarm/resources/types"
)

// Runner to check consecutive if the radio should be started or not
func Runner() {
	for {
		fmt.Println("Infinite Loop")
		time.Sleep(time.Second * 5)
		// TODO: implement logic
	}
}

// SetUpRouter set up router for alarm endpoints
func SetUpRouter(router *mux.Router) {
	router.HandleFunc("alarm", getAlarms).Methods(http.MethodGet)
	router.HandleFunc("alarm", addAlarm).Methods(http.MethodPost)

	router.HandleFunc("alarm/{idx}", getAlarm).Methods(http.MethodGet)
	router.HandleFunc("alarm/{idx}", changeAlarm).Methods(http.MethodPut)
	router.HandleFunc("alarm/{idx}", deleteAlarm).Methods(http.MethodDelete)
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
	// TODO: write alarms to file
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
