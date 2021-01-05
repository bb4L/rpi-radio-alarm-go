package alarmtypes

// Alarm struct to represent a alarm
type Alarm struct {
	Name   string `yaml:"name" json:"name"`
	Hour   int    `yaml:"hour" json:"hour"`
	Minute int    `yaml:"minute" json:"minute"`
	Days   []int  `yaml:"days" json:"days"`
	Active bool   `yaml:"active" json:"on"`
}

// Radio struct to represent a radio
type Radio struct {
	Running     bool   `yaml:"running"`
	Pid         int    `yaml:"pid"`
	BodyRunning string `json:"switch"`
}
