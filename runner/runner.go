package runner

import (
	"os"
	"time"

	"github.com/bb4L/rpi-radio-alarm-go-library/storage"
	"github.com/bb4L/rpi-radio-alarm-go-library/types"
	"github.com/bb4L/rpi-radio-alarm-go/constants"

	"github.com/bb4L/rpi-radio-alarm-go-library/logging"
)

var logger = logging.GetLogger(os.Stdout, constants.DefaultPrefix, "runner")

// Runner to check consecutive if the radio should be started or stopped
func Runner(storageHelper *storage.Helper) {
	logger.Println("starting")
	lastShouldRun := false
	for {
		// convert that weekday = 0 is monday and not sunday
		weekdayGo := int(time.Now().Weekday())
		weekday := weekdayGo - 1
		if weekday < 0 {
			weekday = 6
		}

		// storedData, err := storageHelper.GetStoredData(false)

		alarms, err := storageHelper.GetAlarms(false)

		if err != nil {
			logger.Fatalf("Error while loading data: %s", err)
		}

		actualTime := time.Now()
		shouldRun := false

		for _, alarm := range alarms {
			if !alarm.Active {
				continue
			}

			referenceTime := time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day(), alarm.Hour, alarm.Minute, 0, 0, actualTime.Location())
			diff1 := actualTime.Sub(referenceTime)

			referenceTimeNextDay := time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day(), alarm.Hour, alarm.Minute, 0, 0, actualTime.Location()).AddDate(0, 0, 1)
			diff2 := actualTime.Sub(referenceTimeNextDay)

			if checkIfShoulBeRunning(diff1, weekday, alarm) || checkIfShoulBeRunning(diff2, weekday, alarm) {
				shouldRun = true
			}

		}

		if shouldRun != lastShouldRun {
			_, err := storageHelper.SwitchRadio(shouldRun)
			if err != nil {
				logger.Println("error switching radio", err)
			}
			// changeRunning(storageHelper, shouldRun)
			lastShouldRun = shouldRun

		}

		// wait a seconds
		time.Sleep(time.Second)
	}
}

func checkIfShoulBeRunning(diff time.Duration, weekday int, alarm types.Alarm) (shouldRun bool) {
	shouldRun = false
	if diff.Seconds() >= 0 && diff.Minutes() < 5 {
		sameDay := false

		for _, day := range alarm.Days {
			if weekday == day {
				sameDay = true
				break
			}
		}

		if sameDay {
			shouldRun = true
			return
		}
	}
	return
}

// func changeRunning(storageHelper *storage.StorageHelper, shouldRun bool) {

// 	storedData, err := storageHelper.GetStoredData(true)
// 	if err != nil {
// 		logger.Fatalf("Error while loading data: %s", err)
// 	}
// 	if shouldRun {
// 		storedData.Radio.StartRadio()
// 	} else {
// 		err := storedData.Radio.StopRadio()
// 		if err != nil {
// 			logger.Printf("Error while stopping radio %s", err)
// 		}
// 	}

// 	storageHelper.SaveStoredData(storedData)
// }
