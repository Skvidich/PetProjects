package StatusRelay

import (
	"dataCollector/common"
	"sync"
)

type StatusRelay struct {
	inputChan chan common.StatusResponse
	isLog     bool
	muLog     sync.Mutex
	isResend  bool
	muResend  sync.Mutex
	endChan   chan struct{}
}

func NewStatusRelay(input chan common.StatusResponse, isLog bool, isResend bool) *StatusRelay {
	return &StatusRelay{
		inputChan: input,
		isLog:     isLog,
		muLog:     sync.Mutex{},
		isResend:  isResend,
		muResend:  sync.Mutex{},
		endChan:   make(chan struct{}),
	}
}

func (rel *StatusRelay) inputProcess() {
	for status := range rel.inputChan {

		if rel.GetLogState() {
			common.LogStatus(status)
		}

		if rel.GetResendState() {

		}
	}
	rel.endChan <- struct{}{}
}

func (rel *StatusRelay) Close() {
	<-rel.endChan
}
