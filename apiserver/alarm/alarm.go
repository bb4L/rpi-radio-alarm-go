// Package alarm contains all the helping functions to set up the api for the alarm resource
package alarm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bb4L/rpi-radio-alarm-go-library/logging"
	"github.com/bb4L/rpi-radio-alarm-go-library/storage"

	"github.com/bb4L/rpi-radio-alarm-go/apiserver/helper"
	"github.com/bb4L/rpi-radio-alarm-go/constants"

	alarmtypes "github.com/bb4L/rpi-radio-alarm-go-library/types"

	"github.com/gorilla/mux"
)

var logger = logging.GetLogger(os.Stdout, constants.DefaultPrefix, "apiserver", "alarm")

var storageHelper *storage.Helper

// SetUpRouter sets up the router for the alarm endpoints
func SetUpRouter(router *mux.Router, storageH *storage.Helper) {
	logger.Println("setup router")
	storageHelper = storageH
	router.HandleFunc("", getAlarms).Methods(http.MethodGet)
	router.HandleFunc("", addAlarm).Methods(http.MethodPost)

	router.HandleFunc("/{idx}", getAlarm).Methods(http.MethodGet)
	router.HandleFunc("/{idx}", changeAlarm).Methods(http.MethodPut)
	router.HandleFunc("/{idx}", deleteAlarm).Methods(http.MethodDelete)
}

func getAlarms(w http.ResponseWriter, r *http.Request) {
	alarms, err := storageHelper.GetAlarms(false)
	if err != nil {
		helper.HandleStorageError(w, err)
		return
	}

	jsonData, err := json.Marshal(alarms)
	if err != nil {
		helper.HandleJSONMarshalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func addAlarm(w http.ResponseWriter, r *http.Request) {
	var err error
	var alarm alarmtypes.Alarm

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HandleReadAllError(w, err)
		return
	}

	err = json.Unmarshal(body, &alarm)
	if err != nil {
		helper.HandleUnMarshalError(w, err)
		return
	}

	alarms, err := storageHelper.AddAlarm(alarm)
	if err != nil {
		helper.HandleStorageError(w, err)
		return
	}

	jsonData, err := json.Marshal(alarms)
	if err != nil {
		helper.HandleJSONMarshalError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getAlarm(w http.ResponseWriter, r *http.Request) {

	alarms, err := storageHelper.GetAlarms(false)

	if err != nil {
		helper.HandleStorageError(w, err)
		return
	}
	pathParams := mux.Vars(r)

	// TODO: move the parseIdx inside the storage helper
	alarmIdx, errIdx := helper.ParseIdx(w, pathParams, alarms)
	if errIdx != nil {
		logger.Printf("Received error: %s\n", errIdx)
		return
	}

	alarm, err := storageHelper.GetAlarm(alarmIdx, false)
	if err != nil {
		helper.HandleStorageError(w, err)
		return
	}

	jsonData, err := json.Marshal(alarm)
	if err != nil {
		helper.HandleJSONMarshalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func changeAlarm(w http.ResponseWriter, r *http.Request) {

	alarms, err := storageHelper.GetAlarms(false)

	if err != nil {
		helper.HandleStorageError(w, err)
		return
	}

	pathParams := mux.Vars(r)

	alarmIdx, errIdx := helper.ParseIdx(w, pathParams, alarms)
	if errIdx != nil {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HandleReadAllError(w, err)
		return
	}

	var alarm alarmtypes.Alarm
	err = json.Unmarshal(body, &alarm)
	if err != nil {
		helper.HandleUnMarshalError(w, err)
		return
	}

	savedAlarm, err := storageHelper.SaveAlarm(alarmIdx, alarm)
	if err != nil {
		helper.HandleStorageError(w, err)
		return
	}

	jsonData, err := json.Marshal(savedAlarm)
	if err != nil {
		helper.HandleJSONMarshalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteAlarm(w http.ResponseWriter, r *http.Request) {
	alarms, err := storageHelper.GetAlarms(false)

	if err != nil {
		helper.HandleStorageError(w, err)
		return
	}

	pathParams := mux.Vars(r)

	alarmIdx, errIdx := helper.ParseIdx(w, pathParams, alarms)
	if errIdx != nil {
		return
	}

	alarms, err = storageHelper.DeleteAlarm(alarmIdx)
	if err != nil {
		helper.HandleStorageError(w, err)
		return
	}

	jsonData, err := json.Marshal(alarms)
	if err != nil {
		helper.HandleJSONMarshalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
