package storage

import (
	"io/ioutil"

	"github.com/bb4L/rpi-radio-alarm-go/logging"

	alarmtypes "github.com/bb4L/rpi-radio-alarm-go-library/types"

	"gopkg.in/yaml.v2"
)

// Settings contains specific settings for the program
type Settings struct {
	Port int `yaml:"port"`
}

// RpiRadioAlarmData contains all the data for the RpiRadioAlarm
type RpiRadioAlarmData struct {
	Settings Settings           `yaml:"settings"`
	Alarms   []alarmtypes.Alarm `yaml:"alarms"`
	Radio    alarmtypes.Radio   `yaml:"radio"`
}

const dataFilename = "./rpi_data.yaml"

// Returns the whole data
func GetStoredData() (RpiRadioAlarmData, error) {
	fileData, err := ioutil.ReadFile(dataFilename)

	if err != nil {
		panic(err)
	}

	var data RpiRadioAlarmData

	source := []byte(fileData)
	err = yaml.Unmarshal(source, &data)
	if err != nil {
		logging.GetFatalLogger().Fatalf("error: %v", err)
	}
	return data, err
}

// Save the data to the file
func SaveData(data RpiRadioAlarmData) {
	outSource, err := yaml.Marshal(data)
	if err != nil {
		logging.GetFatalLogger().Fatalf("error: %v", err)
	}

	ioutil.WriteFile(dataFilename, outSource, 0777)
}
