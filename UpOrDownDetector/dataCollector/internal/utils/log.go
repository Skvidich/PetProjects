package utils

import (
	"dataCollector/pkg/types"
	"log"
	"os"
	"time"
)

var StatLogger *log.Logger
var ErrLogger *log.Logger
var StatLogDone = make(chan struct{})
var ErrLogDone = make(chan struct{})

func InitStatLogger(path string) {
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	StatLogger = log.New(file, "", log.Ldate|log.Ltime)
	go func() {
		<-StatLogDone
		file.Close()
	}()
}

func InitErrLog(path string) {
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	ErrLogger = log.New(file, "", log.Ldate|log.Ltime)
	go func() {
		<-ErrLogDone
		file.Close()
	}()
}

func LogStatus(status types.StatusResponse) {
	StatLogger.Println(status.Name, " : ", status.Time.String())
	for _, component := range status.Components {
		StatLogger.Println("\t", component.Name, " : ", component.Status)
	}
	StatLogger.Println()
}

func LogError(errStr string) {
	ErrLogger.Println(errStr, " : ", time.Now().String())
	ErrLogger.Println()
}

func KillStatLogger() {
	StatLogDone <- struct{}{}
}

func KillErrLogger() {
	ErrLogDone <- struct{}{}
}
