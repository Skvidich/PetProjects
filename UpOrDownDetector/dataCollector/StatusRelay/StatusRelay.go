package StatusRelay

import (
	"dataCollector/common"
)

type StatusRelay struct {
	inputChan     chan common.StatusResponse
	isLog         bool
	isResend      bool
	endChan       chan struct{}
	processStatus StatusProcess
}

func NewStatusRelay(input chan common.StatusResponse, isLog bool, isResend bool) *StatusRelay {

	processStatus := func(response common.StatusResponse) {}

	if isLog {
		logWrap(processStatus)
	}

	if isResend {
		resendWrap(processStatus)
	}
	return &StatusRelay{
		inputChan:     input,
		isLog:         isLog,
		isResend:      isResend,
		endChan:       make(chan struct{}),
		processStatus: processStatus,
	}
}

func (rel *StatusRelay) inputProcess() {
	for status := range rel.inputChan {

		rel.processStatus(status)
	}
	rel.endChan <- struct{}{}
}

type StatusProcess func(st common.StatusResponse)

func logWrap(prevFunc StatusProcess) StatusProcess {
	return func(status common.StatusResponse) {
		common.LogStatus(status)
		prevFunc(status)
	}
}

func resendWrap(prevFunc StatusProcess) StatusProcess {
	return func(status common.StatusResponse) {
		// ResendStatusFunction
		prevFunc(status)
	}
}
