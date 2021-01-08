package runner

import (
	storage "rpi-radio-alarm/helper"
	"rpi-radio-alarm/logging"
	"time"
)

// Runner to check consecutive if the radio should be started or not
func Runner() {
	logging.GetInfoLogger().Println("start runner...")
	for {
		// convert that weekday = 0 is monday and not sunday
		weekdayGo := int(time.Now().Weekday())
		weekday := weekdayGo - 1
		if weekday < 0 {
			weekday = 6
		}

		storedData, err := storage.GetStoredData()
		if err != nil {
			logging.GetFatalLogger().Printf("Error while loading data: %s", err)
		}

		actualTime := time.Now()
		shouldRun := false

		// var secondsDiff float64

		for _, alarm := range storedData.Alarms {
			if !alarm.Active {
				continue
			}

			referenceTime := time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day(), alarm.Hour, alarm.Minute, 0, 0, actualTime.Location())
			diff1 := actualTime.Sub(referenceTime)

			if diff1.Seconds() >= 0 && diff1.Minutes() < 5 {
				sameDay := false

				for _, day := range alarm.Days {
					if weekday == day {
						sameDay = true
						break
					}
				}

				if sameDay {
					shouldRun = true
					// secondsDiff = diff1.Seconds()
					break
				}
			}

			referenceTimeNextDay := time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day(), alarm.Hour, alarm.Minute, 0, 0, actualTime.Location()).AddDate(0, 0, 1)
			diff2 := actualTime.Sub(referenceTimeNextDay)

			if diff2.Seconds() >= 0 && diff2.Minutes() < 5 {
				sameDay := false

				for _, day := range alarm.Days {
					if weekday == day {
						sameDay = true
						break
					}
				}

				if sameDay {
					shouldRun = true
					// secondsDiff = diff2.Seconds()
					break
				}
			}

		}

		logging.GetInfoLogger().Println("Alarms checked")

		if shouldRun != storedData.Radio.Running {
			if shouldRun {
				storedData.Radio.StartRadio()

				// TODO: implement a good way to ensure the radio running for 5min if it's a alarm and endless if it's started manually

				// } else {
				// 	err := storedData.Radio.StopRadio()
				// 	if err != nil {
				// 		logging.GetFatalLogger().Printf("Error while stopping radio %s", err)
				// 	}
			}
			storage.SaveData(storedData)
		}

		// wait a seconds
		time.Sleep(time.Second)
	}
}
