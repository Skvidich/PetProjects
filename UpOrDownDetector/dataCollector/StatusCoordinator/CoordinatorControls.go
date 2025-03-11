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

	return GetterInfo{Name: name, State: cord.getterList[name].GetState().String()}
}

func (cord *StatusCoordinator) GetListInfo() []GetterInfo {

	res := make([]GetterInfo, len(cord.getterList))
	for name, getter := range cord.getterList {
		res = append(res, GetterInfo{Name: name, State: getter.GetState().String()})
	}

	return res

}

func (cord *StatusCoordinator) RunGetter(name string) {
	cord.addGetter(name)
	cord.getterList[name].SetState(StatusGetters.IsRunning)
	go cord.getterList[name].RunProcess()
}

func (cord *StatusCoordinator) RunAll() {
	for _, name := range cord.getterNames {
		cord.addGetter(name)
		cord.getterList[name].SetState(StatusGetters.IsRunning)
		go cord.getterList[name].RunProcess()
	}
}

func (cord *StatusCoordinator) StopAll() {
	for _, getter := range cord.getterList {
		getter.SetState(StatusGetters.IsDown)
	}
}

func (cord *StatusCoordinator) StopGetter(name string) {
	cord.getterList[name].SetState(StatusGetters.IsDown)
}

func (cord *StatusCoordinator) Shutdown() {
	for _, getter := range cord.getterList {
		getter.SetState(StatusGetters.IsDown)
	}

	for cord.getterCount() != 0 {
		// maybe add time.Sleep() there
	}
	close(cord.OutChan)
	close(cord.feedbackChan)
}

func (cord *StatusCoordinator) IsGetterExist(name string) bool {
	cord.muList.Lock()
	defer cord.muList.Unlock()
	_, res := cord.getterList[name]
	return res
}
