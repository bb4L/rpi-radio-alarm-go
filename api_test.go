package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"

	alarmtypes "github.com/bb4L/rpi-radio-alarm-go-library/types"

	"testing"

	"github.com/bb4L/rpi-radio-alarm-go/resources/alarm"
	"github.com/bb4L/rpi-radio-alarm-go/resources/radio"

	"github.com/gorilla/mux"
)

func TestApi(t *testing.T) {
	// Create server using the a router initialized elsewhere. The router
	// can be a Gorilla mux as in the question, a net/http ServeMux,
	// http.DefaultServeMux or any value that statisfies the net/http
	// Handler interface.
	r := mux.NewRouter()

	alarm.SetUpRouter(r.PathPrefix("/alarm").Subrouter())
	radio.SetUpRouter(r.PathPrefix("/radio").Subrouter())
	ts := httptest.NewServer(r)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	alarm1 := alarmtypes.Alarm{Name: "Test", Hour: 7, Minute: 0, Active: false, Days: []int{0, 1}}
	alarm2 := alarmtypes.Alarm{Name: "Test2", Hour: 9, Minute: 15, Active: false, Days: []int{0, 1}}
	alarm3 := alarmtypes.Alarm{Name: "Test3", Hour: 11, Minute: 5, Active: true, Days: []int{0, 1}}
	alarm3Json, _ := json.Marshal(alarm3)

	alarm3Changed := alarm3
	alarm3Changed.Name = "Test3.2"
	alarm3ChangedJson, _ := json.Marshal(alarm3Changed)

	result1, _ := json.Marshal([]alarmtypes.Alarm{alarm1, alarm2})
	resultAddAlarm, _ := json.Marshal([]alarmtypes.Alarm{alarm1, alarm2, alarm3})
	radio := alarmtypes.Radio{Pid: -1, Running: false}
	result2, _ := json.Marshal(radio)

	// startRadio, _ := json.Marshal(map[string]string{"switch": "on"})
	stopRadio, _ := json.Marshal(map[string]string{"switch": "off"})
	stopRadioResult, _ := json.Marshal(alarmtypes.Radio{Pid: -1, Running: false})

	tests := []struct {
		name          string
		r             *http.Request
		expected_data []byte
	}{
		{name: "get alarm", r: newreq("GET", ts.URL+"/alarm", nil), expected_data: result1},
		{name: "get radio", r: newreq("GET", ts.URL+"/radio", nil), expected_data: result2},
		{name: "add alarm", r: newreq("POST", ts.URL+"/alarm", bytes.NewReader(alarm3Json)), expected_data: resultAddAlarm},
		{name: "change alarm", r: newreq("PUT", ts.URL+"/alarm/2", bytes.NewReader(alarm3ChangedJson)), expected_data: alarm3ChangedJson},

		{name: "delete alarm", r: newreq("DELETE", ts.URL+"/alarm/2", nil), expected_data: result1},
		// TODO: mock call to radio running
		// {name: "start radio", r: newreq("POST", ts.URL+"/radio", bytes.NewReader(startRadio)), expected_data: stopRadioResult},
		{name: "stop radio", r: newreq("POST", ts.URL+"/radio", bytes.NewReader(stopRadio)), expected_data: stopRadioResult},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.r)
			if err != nil {
				t.Fatal(err)
			}

			body, _ := ioutil.ReadAll(resp.Body)

			defer resp.Body.Close()
			var data interface{}
			json.Unmarshal(body, &data)
			if err != nil {
				panic(err)
			}

			if !reflect.DeepEqual(body, tt.expected_data) {
				t.Fatalf("recieved %s instead of %s", body, tt.expected_data)
			}
		})
	}
}
