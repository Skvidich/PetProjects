package StatusRelay

import (
	"dataCollector/common"
	"fmt"
	"sync"
)

type StatusRelay struct {
	inputChan chan common.StatusResponse
	isLog     bool
	muLog     sync.Mutex
	isResend  bool
	muResend  sync.Mutex
	logPath   string
}

func NewStatusRelay(input chan common.StatusResponse, isLog bool, isResend bool, logPath string) *StatusRelay {
	return &StatusRelay{
		inputChan: input,
		isLog:     isLog,
		muLog:     sync.Mutex{},
		isResend:  isResend,
		muResend:  sync.Mutex{},
		logPath:   logPath,
	}
}

func (rel *StatusRelay) inputProcess() {
	for status := range rel.inputChan {

		if rel.GetLogState() {
			logStatus(status)
		}

		if rel.GetResendState() {

		}
	}
}

func logStatus(status common.StatusResponse) {
	fmt.Println(status.Name, " : ", status.Time.String())
	for _, component := range status.Components {
		fmt.Println(component.Name, " : ", component.Status)
	}
	fmt.Println()
}
