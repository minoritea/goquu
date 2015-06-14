package goquu

import (
	loglib "log"
	"os"
)

var logger *loglib.Logger = loglib.New(os.Stderr, "Log: ", 0)

func InitLogger(new_logger *loglib.Logger) {
	logger = new_logger
}

func SetLoggerFromFile(path string, logflags int) (err error) {
	file, err := os.OpenFile(path,  os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0776)
	if err != nil {
		return
	}
	_, err = file.WriteString("\nStart logging...\n\n")
	if err != nil {
		file.Close()
		return
	}
	InitLogger(loglib.New(file, "Log: ", logflags))
	return
}
