package goquu

import (
	loglib "log"
	"os"
)

var logger *loglib.Logger = loglib.New(os.Stderr, "Log: ", 0)

func InitLogger(new_logger *loglib.Logger) {
	logger = new_logger
}
