package common

import (
	"log"
	"os"
	"time"
)

var StatusLogger *log.Logger
var ErrorLogger *log.Logger
var StatusLoggerChan = make(chan struct{})
var StatusErrorChan = make(chan struct{})

func StatusLoggerInit(path string) {
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	StatusLogger = log.New(file, "", log.Ldate|log.Ltime)
	go func() {
		<-StatusLoggerChan
		file.Close()
	}()
}

func ErrorLoggerInit(path string) {
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	ErrorLogger = log.New(file, "", log.Ldate|log.Ltime)
	go func() {
		<-StatusErrorChan
		file.Close()
	}()
}

func LogStatus(status StatusResponse) {
	StatusLogger.Println(status.Name, " : ", status.Time.String())
	for _, component := range status.Components {
		StatusLogger.Println("\t", component.Name, " : ", component.Status)
	}
	StatusLogger.Println()
}

func LogError(errStr string) {
	ErrorLogger.Println(errStr, " : ", time.Now().String())
	ErrorLogger.Println()
}

func KillStatusLogger() {
	StatusLoggerChan <- struct{}{}
}

func KillErrorLogger() {
	StatusErrorChan <- struct{}{}
}
