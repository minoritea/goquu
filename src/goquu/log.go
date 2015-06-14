package goquu
import (
	loglib "log"
	"os"
)

var logger *loglib.Logger = loglib.New(os.Stderr, "Log: ",  loglib.Llongfile)

func InitLogger(new_logger *loglib.Logger) {
	logger = new_logger
}
