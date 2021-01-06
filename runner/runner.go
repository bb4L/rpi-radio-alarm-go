package runner

import (
	"rpi-radio-alarm/logging"
	"time"
)

// Runner to check consecutive if the radio should be started or not
func Runner() {
	logging.GetInfoLogger().Println("start runner...")

	for {
		logging.GetInfoLogger().Println("Infinite Loop")
		// fmt.Println(time.Now())
		time.Sleep(time.Second * 5)
		// TODO: implement logic
	}
}
