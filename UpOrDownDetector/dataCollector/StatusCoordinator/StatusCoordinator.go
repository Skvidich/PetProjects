package StatusCoordinator

import (
	"dataCollector/StatusCoordinator/StatusGetters"
	"dataCollector/common"
	"fmt"
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
	GetterList   map[string]*StatusGetters.StatusGetter
	muList       sync.Mutex
	GetterNames  []string
	OutChan      chan common.StatusResponse
	FeedbackChan chan StatusGetters.GetterFeedback
}

func NewStatusCoordinator() *StatusCoordinator {

	outChan := make(chan common.StatusResponse)

	getterNames := iniGetterNames()
	feedbackChan := make(chan StatusGetters.GetterFeedback, len(getterNames))

	getterList := make(map[string]*StatusGetters.StatusGetter, len(getterNames))

	for _, name := range getterNames {
		getterList[name] = StatusGetters.NewStatusGetter(name, StatusGetters.GetterFuncList[name], time.Second*300, &feedbackChan)
		attachToBus(outChan, getterList[name].OutputChan)
	}

	res := &StatusCoordinator{
		GetterNames:  getterNames,
		muList:       sync.Mutex{},
		GetterList:   getterList,
		OutChan:      outChan,
		FeedbackChan: feedbackChan,
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

	for feedback := range cord.FeedbackChan {

		if feedback.Err != nil {
			logGetterError(feedback)
		} else {
			cord.removeGetter(feedback.Name)
		}

		//cord.GetterList[feedback.Name].SetState(StatusGetters.IsRunning)
		//go cord.GetterList[feedback.Name].RunProcess()

	}
}

func logGetterError(getErr StatusGetters.GetterFeedback) {
	fmt.Println(getErr.Name, " ", getErr.Err.Error())
}

func iniGetterNames() []string {
	res := make([]string, 0)
	res = append(res, "Github")
	return res
}

func (cord *StatusCoordinator) removeGetter(name string) {
	cord.muList.Lock()
	defer cord.muList.Unlock()
	delete(cord.GetterList, name)

}

func (cord *StatusCoordinator) addGetter(name string) {
	cord.muList.Lock()
	defer cord.muList.Unlock()
	cord.GetterList[name] = StatusGetters.NewStatusGetter("name", StatusGetters.GetterFuncList[name], time.Second*300, &cord.FeedbackChan)

}

func (cord *StatusCoordinator) getterCount() int {
	cord.muList.Lock()
	defer cord.muList.Unlock()
	return len(cord.GetterList)

}
