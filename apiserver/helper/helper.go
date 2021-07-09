// Package helper contains some helper functions for the APIs
package helper

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bb4L/rpi-radio-alarm-go-library/logging"
	"github.com/bb4L/rpi-radio-alarm-go-library/types"
	"github.com/bb4L/rpi-radio-alarm-go/constants"
)

var logger = logging.GetLogger(os.Stdout, constants.DefaultPrefix, "apiserver", "helper")

// ParseIdx parses the idx from the given params
func ParseIdx(w http.ResponseWriter, pathParams map[string]string, alarms []types.Alarm) (alarmIdx int, err error) {
	alarmIdx = -1

	if val, ok := pathParams["idx"]; ok {
		alarmIdx, err = strconv.Atoi(val)
		if err != nil {
			parseIdxError(w, val, err)
			return
		}
	}

	if isIndexInvalid(alarms, alarmIdx) {
		err = fmt.Errorf("index not valid %d", alarmIdx)
		handleInvalidIdx(w, err)
		return
	}

	return
}

func isIndexInvalid(alarms []types.Alarm, alarmIdx int) bool {
	return alarmIdx < 0 || alarmIdx >= len(alarms)
}

func handleInvalidIdx(w http.ResponseWriter, errIdx error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, errIdx.Error())))
}

// HandleStorageError handles the given error with the given ResponseWriter
func HandleStorageError(w http.ResponseWriter, err error) {

	logger.Printf("storage error %s", err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "error on getting stored data"}`))
}

// HandleJSONMarshalError handles the given error with the given ResponseWriter
func HandleJSONMarshalError(w http.ResponseWriter, err error) {
	logger.Printf("json error %s", err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "error on creating json data"}`))
}

// HandleUnMarshalError handles the given error with the given ResponseWriter
func HandleUnMarshalError(w http.ResponseWriter, err error) {
	logger.Printf("error unmarshalling body %s", err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "error on unmarshal body"}`))
}

// HandleReadAllError handles the given error with the given ResponseWriter
func HandleReadAllError(w http.ResponseWriter, err error) {
	logger.Printf("error while reading body %s", err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "error on reading body"}`))
}

func parseIdxError(w http.ResponseWriter, val string, err error) {
	logger.Printf("error parsing index %s %s", val, err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "need a number as path param"}`))
}
