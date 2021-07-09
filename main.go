package main

import (
	"os"
	"strconv"

	"github.com/bb4L/rpi-radio-alarm-go-library/api"
	"github.com/bb4L/rpi-radio-alarm-go-library/storage"
	server "github.com/bb4L/rpi-radio-alarm-go/apiserver"
	"github.com/bb4L/rpi-radio-alarm-go/constants"

	"github.com/bb4L/rpi-radio-alarm-go-library/logging"
	"github.com/bb4L/rpi-radio-alarm-go/runner"
	bot "github.com/bb4L/rpi-radio-alarm-telegrambot-go/bot"
)

var logger = logging.GetLogger(os.Stdout, constants.DefaultPrefix, "main")
var storageHelper storage.Helper = storage.Helper{}

func main() {

	settings, err := storageHelper.GetSettings()

	if err != nil {
		logger.Fatalln("error on getting stored data")
	}

	if settings.RunAPI {
		logger.Println("start api server")
		go server.StartAPIServer(&storageHelper)
	}

	if settings.RunTelegrambot {
		apiHelper := api.Helper{AlarmURL: "http://localhost:" + strconv.Itoa(settings.Port), ExtraHeader: "", ExtreaHeaderValue: ""}
		go bot.StartTelegramBot(&apiHelper)
	}

	if settings.RunDiscordbot {
		go func() {
			logger.Println("run dicordbot")
		}()
	}

	runner.Runner(&storageHelper)
}
