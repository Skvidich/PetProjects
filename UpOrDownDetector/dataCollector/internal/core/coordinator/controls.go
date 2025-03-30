package coordinator

import (
	"dataCollector/internal/getters"
)

type GetterInfo struct {
	Name  string
	State string
}

func (c *Coordinator) Getter(name string) GetterInfo {

	return GetterInfo{Name: name, State: c.getters[name].State().String()}
}

func (c *Coordinator) GetterList() []GetterInfo {

	res := make([]GetterInfo, len(c.getters))
	for name, getter := range c.getters {
		res = append(res, GetterInfo{Name: name, State: getter.State().String()})
	}

	return res

}

func (c *Coordinator) Start(name string) {
	c.addGetter(name)
	c.getters[name].SetState(getters.Running)
	go c.getters[name].Start()
}

func (c *Coordinator) StartAll() {
	for _, name := range c.getterNames {
		c.addGetter(name)
		c.getters[name].SetState(getters.Running)
		go c.getters[name].Start()
	}
}

func (c *Coordinator) StopAll() {
	for _, getter := range c.getters {
		getter.SetState(getters.Down)
	}
}

func (c *Coordinator) Stop(name string) {
	c.getters[name].SetState(getters.Down)
}

func (c *Coordinator) Shutdown() {
	for _, getter := range c.getters {
		getter.SetState(getters.Down)
	}

	for c.count() != 0 {
		// maybe add time.Sleep() there
	}
	close(c.OutChan)
	close(c.feedbackChan)
}

func (c *Coordinator) Exists(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, getter := range c.getterNames {
		if getter == name {
			return true
		}
	}
	return false
}
