package alarmtypes

import (
	"os"
	"os/exec"
	"rpi-radio-alarm/logging"
)

// Alarm struct to represent a alarm
type Alarm struct {
	Name   string `yaml:"name" json:"name"`
	Hour   int    `yaml:"hour" json:"hour"`
	Minute int    `yaml:"minute" json:"min"`
	Days   []int  `yaml:"days" json:"days"`
	Active bool   `yaml:"active" json:"on"`
}

// Radio struct to represent a radio
type Radio struct {
	Running bool `yaml:"running" json:"isPlaying"`
	Pid     int  `yaml:"pid"`
}

// StartRadio start the radio
func (r *Radio) StartRadio() {
	// TODO: exec command to start radio
	cmd := exec.Command("while true; do (echo \"test\" &&  sleep 5); done")
	// cmd := exec.Command("mplayer https://streamingp.shoutcast.com/hotmixradio-sunny-128.mp3 -volume 150")
	cmd.Stdout = logging.GetInfoLogger().Writer()
	cmd.Stderr = logging.GetErrorLogger().Writer()

	r.Pid = cmd.Process.Pid
	r.Running = true
}

// StopRadio stop the radio
func (r *Radio) StopRadio() error {
	var err error

	if r.Pid == -1 {
		return nil
	}

	process, err := os.FindProcess(r.Pid)
	if err != nil {
		return err
	}

	err = process.Kill()

	if err != nil {
		return err
	}

	r.Running = false
	r.Pid = -1

	return nil
}
