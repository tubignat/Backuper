package logging

import (
	"io"
	"log"
	"os"
)

func Configure(logfilepath string) {
	f, err := os.OpenFile(logfilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error("Could not open a log file", err)
		return
	}
	writer := io.MultiWriter(f, os.Stderr)
	log.SetOutput(writer)
}

func Debug(message ...interface{}) {
	log.Print("DEBUG    ", message)
}

func Info(message ...interface{}) {
	log.Print("INFO     ", message)
}

func Warning(message ...interface{}) {
	log.Print("WARNING  ", message)
}

func Error(message ...interface{}) {
	log.Print("ERROR    ", message)
}
