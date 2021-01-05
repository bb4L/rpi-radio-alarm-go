package radio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"rpi-radio-alarm/logging"
	alarmtypes "rpi-radio-alarm/resources/types"

	"github.com/gorilla/mux"
)

// GetRadio returns the radio
func GetRadio() alarmtypes.Radio {
	// TODO: read radio state from file
	return alarmtypes.Radio{}
}

func postRadio(w http.ResponseWriter, r *http.Request) {
	logging.GetInfoLogger().Println(r)

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var msg alarmtypes.Radio
	err = json.Unmarshal(b, &msg)
	logging.GetInfoLogger().Print(b)
	logging.GetInfoLogger().Print(err)
	logging.GetInfoLogger().Print(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// TODO: handle value from body depending on actual value

	logging.GetInfoLogger().Print(r.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// SetUpRouter set up router for radio endpoints
func SetUpRouter(r *mux.Router) {
	r.HandleFunc("/", postRadio).Methods(http.MethodPost)
}

func startRadio() {
	cmd := exec.Command("while true; do (echo \"test\" &&  sleep 5); done")
	cmd.Stdout = logging.GetInfoLogger().Writer()
	cmd.Stderr = logging.GetErrorLogger().Writer()
	err := cmd.Start()

	// cmd.Process.Pid
	if err != nil {
		logging.GetFatalLogger().Printf("failed with: %s", err)
	}
	// TODO: write pid to file?
	// TODO: check if file with pid exists
}
