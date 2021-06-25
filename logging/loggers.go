// Package including some logging utilities
package logging

import (
	"log"
	"os"
	"sync"
)

const commonPrefix = "[rpi-radio-alarm-go] "

var (
	// InfoLogger logging info
	InfoLogger *log.Logger

	// ErrorLogger logging errors
	ErrorLogger *log.Logger

	// FatalLogger logging fatal
	FatalLogger *log.Logger

	once sync.Once
)

// Get the InfoLogger
func GetInfoLogger() *log.Logger {
	once.Do(func() {
		InfoLogger = log.New(os.Stderr, commonPrefix+"INFO: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	})
	return InfoLogger
}

// Get the ErrorLogger
func GetErrorLogger() *log.Logger {
	once.Do(func() {
		ErrorLogger = log.New(os.Stderr, commonPrefix+"ERROR: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	})
	return ErrorLogger
}

// Get the FatalLogger
func GetFatalLogger() *log.Logger {
	once.Do(func() {
		FatalLogger = log.New(os.Stderr, commonPrefix+"FATAL: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	})
	return FatalLogger
}
