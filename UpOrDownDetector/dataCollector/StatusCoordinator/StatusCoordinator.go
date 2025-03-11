package StatusCoordinator

import (
	"dataCollector/StatusCoordinator/StatusGetters"
	"time"
)

// TO-DO
// Add initialising Status coordinator by ini file
//

type StatusCoordinator struct {
	GetterList map[string]*StatusGetters.StatusGetter
	OutChan    chan StatusGetters.StatusResponse
	ErrorChan  chan StatusGetters.GetterError
}

func NewStatusCoordinator() StatusCoordinator {

	outChan := make(chan StatusGetters.StatusResponse)

	getterList := make(map[string]*StatusGetters.StatusGetter, 0)
	errorChan := make(chan StatusGetters.GetterError)
	getterList["Github"] = StatusGetters.NewStatusGetter("Github", StatusGetters.GetterFuncList["Github"], time.Second*300, &errorChan)

	for i := range getterList {
		AttachToBus(outChan, getterList[i].OutputChan)
	}

	return StatusCoordinator{
		GetterList: getterList,
		OutChan:    outChan,
		ErrorChan:  errorChan,
	}
}

func AttachToBus(bus chan StatusGetters.StatusResponse, newChan chan StatusGetters.StatusResponse) {

	go func() {
		for resp := range newChan {
			bus <- resp
		}
	}()

}
