package StatusCoordinator

import "dataCollector/StatusCoordinator/StatusGetters"

type GetterInfo struct {
	Name  string
	State StatusGetters.GetterState
}

func (cord *StatusCoordinator) GetInfo(name string) GetterInfo {
	return GetterInfo{Name: name, State: cord.GetterList[name].GetState()}
}

func (cord *StatusCoordinator) GetListInfo() []GetterInfo {

	res := make([]GetterInfo, len(cord.GetterList))
	for name, getter := range cord.GetterList {
		res = append(res, GetterInfo{Name: name, State: getter.GetState()})
	}

	return res

}

func (cord *StatusCoordinator) RunGetter(name string) {
	cord.GetterList[name].SetState(StatusGetters.IsRunning)
	go cord.GetterList[name].RunProcess()
}
