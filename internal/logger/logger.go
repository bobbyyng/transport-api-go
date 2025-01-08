package logger

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	Logger *log.Logger
)

func Init(logDir string) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	currentDate := time.Now().Format("2006-01-02")
	logFile := logDir + "/log-" + currentDate + ".log"

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	Logger = log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func Success(message string) {
	Logger.Println("[SUCCESS]", message)
}

func Warning(message string) {
	Logger.Println("[WARNING]", message)
}

func Error(message string) {
	Logger.Println("[ERROR]", message)
}

func Fatal(message string) {
	Logger.Fatal("[FATAL]", message)
}
