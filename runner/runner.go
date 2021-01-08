package runner

import (
	"fmt"
	storage "rpi-radio-alarm/helper"
	"rpi-radio-alarm/logging"
	"time"
)

// Runner to check consecutive if the radio should be started or not
func Runner() {
	logging.GetInfoLogger().Println("start runner...")

	for {
		logging.GetInfoLogger().Println("Infinite Loop")
		fmt.Println(time.Now())
		fmt.Println(time.Now().Hour())
		fmt.Println(time.Now().Minute())

		fmt.Println(int(time.Now().Weekday()))

		weekday := int(time.Now().Weekday())

		storedData, err := storage.GetStoredData()
		if err != nil {
			logging.GetFatalLogger().Printf("Error while loading data: %s", err)
		}

		actualTime := time.Now()
		shouldRun := false

		for _, alarm := range storedData.Alarms {
			referenceTime := time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day(), alarm.Hour, alarm.Minute, 0, 0, actualTime.Location())
			referenceTimeNextDay := time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day()+1, alarm.Hour, alarm.Minute, 0, 0, actualTime.Location())

			diff1 := actualTime.Sub(referenceTime)
			diff2 := actualTime.Sub(referenceTimeNextDay)

			if diff1.Seconds() > -4 && diff1.Minutes() < 5 {
				sameDay := false

				for _, day := range alarm.Days {
					if weekday == day {
						sameDay = true
						break
					}
				}

				if sameDay {
					shouldRun = true
					break
				}
			}

			if diff2.Seconds() > -4 && diff2.Minutes() < 5 {
				sameDay := false

				for _, day := range alarm.Days {
					if weekday == day {
						sameDay = true
						break
					}
				}

				if sameDay {
					shouldRun = true
					break
				}
			}

		}

		if shouldRun != storedData.Radio.Running {
			if shouldRun {
				storedData.Radio.StartRadio()
			} else {
				storedData.Radio.StopRadio()
			}
			storage.SaveData(storedData)
		}

		// wait 5 seconds
		time.Sleep(time.Second * 5)

	}
}
