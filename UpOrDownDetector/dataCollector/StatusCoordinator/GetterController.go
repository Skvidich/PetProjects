package StatusCoordinator

import (
	"dataCollector/StatusCoordinator/StatusGetters"
)

type GetterInfo struct {
	Name  string
	State string
}

// TO-DO
// - write new shutdown logic

func (cord *StatusCoordinator) GetInfo(name string) GetterInfo {

	return GetterInfo{Name: name, State: cord.GetterList[name].GetState().String()}
}

func (cord *StatusCoordinator) GetListInfo() []GetterInfo {

	res := make([]GetterInfo, len(cord.GetterList))
	for name, getter := range cord.GetterList {
		res = append(res, GetterInfo{Name: name, State: getter.GetState().String()})
	}

	return res

}

func (cord *StatusCoordinator) RunGetter(name string) {
	cord.addGetter(name)
	cord.GetterList[name].SetState(StatusGetters.IsRunning)
	go cord.GetterList[name].RunProcess()
}

func (cord *StatusCoordinator) RunAll() {
	for _, name := range cord.GetterNames {
		cord.addGetter(name)
		cord.GetterList[name].SetState(StatusGetters.IsRunning)
		go cord.GetterList[name].RunProcess()
	}
}

func (cord *StatusCoordinator) StopAll() {
	for _, getter := range cord.GetterList {
		getter.SetState(StatusGetters.IsDown)
	}
}

func (cord *StatusCoordinator) StopGetter(name string) {
	cord.GetterList[name].SetState(StatusGetters.IsDown)
}

func (cord *StatusCoordinator) Shutdown() {
	for _, getter := range cord.GetterList {
		getter.SetState(StatusGetters.IsDown)
	}

	for cord.getterCount() != 0 {
		// maybe add time.Sleep() there
	}
	close(cord.OutChan)
	close(cord.FeedbackChan)
}
