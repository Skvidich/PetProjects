package StatusCoordinator

import (
	"dataCollector/StatusCoordinator/StatusGetters"
	"dataCollector/common"
	"sync"
	"time"
)

type StatusResponse struct {
	Name       string      `json:"name"`
	Time       time.Time   `json:"time"`
	Components []Component `json:"components"`
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// TO-DO
// Add initialising Status coordinator by ini file
//

type StatusCoordinator struct {
	getterList   map[string]*StatusGetters.StatusGetter
	muList       sync.Mutex
	getterNames  []string
	OutChan      chan common.StatusResponse
	feedbackChan chan StatusGetters.GetterFeedback
	delay        time.Duration
}

func NewStatusCoordinator(delay time.Duration, names []string) *StatusCoordinator {

	outChan := make(chan common.StatusResponse)

	getterNames := names
	feedbackChan := make(chan StatusGetters.GetterFeedback, len(getterNames))

	getterList := make(map[string]*StatusGetters.StatusGetter, len(getterNames))

	res := &StatusCoordinator{
		getterNames:  getterNames,
		muList:       sync.Mutex{},
		getterList:   getterList,
		OutChan:      outChan,
		feedbackChan: feedbackChan,
		delay:        delay,
	}
	go res.processFeedback()
	return res
}

func attachToBus(bus chan common.StatusResponse, newChan chan common.StatusResponse) {

	go func() {
		for resp := range newChan {
			bus <- resp
		}
	}()

}

func (cord *StatusCoordinator) processFeedback() {

	for feedback := range cord.feedbackChan {

		if feedback.Err != nil {
			logGetterError(feedback)
		} else {
			cord.removeGetter(feedback.Name)
		}

		//cord.getterList[feedback.Name].SetState(StatusGetters.IsRunning)
		//go cord.getterList[feedback.Name].RunProcess()

	}
}

func logGetterError(getErr StatusGetters.GetterFeedback) {
	common.LogError(getErr.Name + " " + getErr.Err.Error())
}

func (cord *StatusCoordinator) removeGetter(name string) {
	cord.muList.Lock()
	defer cord.muList.Unlock()
	delete(cord.getterList, name)

}

func (cord *StatusCoordinator) addGetter(name string) {
	cord.muList.Lock()
	defer cord.muList.Unlock()
	cord.getterList[name] = StatusGetters.NewStatusGetter(name, StatusGetters.GetterFuncList[name], cord.delay, &cord.feedbackChan)
	attachToBus(cord.OutChan, cord.getterList[name].OutputChan)
}

func (cord *StatusCoordinator) getterCount() int {
	cord.muList.Lock()
	defer cord.muList.Unlock()
	return len(cord.getterList)

}
