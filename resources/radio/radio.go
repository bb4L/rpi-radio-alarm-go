package radio

import "github.com/gorilla/mux"

// Radio struct to represent a radio
type Radio struct {
	Running bool `yaml:"running"`
}

// GetRadio returns the radio
func GetRadio() Radio {
	return Radio{}
}

// SetUpRouter set up router for radio endpoints
func SetUpRouter(*mux.Router) {

}
