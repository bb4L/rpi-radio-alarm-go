package alarm

import (
	"fmt"
	"time"

	"github.com/gorilla/mux"
)

// Alarm struct to represent a alarm
type Alarm struct {
	Hour   int   `yaml:"hour"`
	Minute int   `yaml:"minute"`
	Days   []int `yaml:"days"`
	Active bool  `yaml:"active"`
}

// Runner to check consecutive if the radio should be started or not
func Runner() {
	for {
		fmt.Println("Infinite Loop")
		time.Sleep(time.Second * 5)
	}
}

// SetUpRouter set up router for alarm endpoints
func SetUpRouter(*mux.Router) {

}

func getAlarms() {

}

func addAlarm() {

}

func deleteAlarm() {

}

func changeAlarm(alarm Alarm) {

}
