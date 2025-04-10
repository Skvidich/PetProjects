package coordinator

import (
	"dataCollector/internal/core/getters"
	"dataCollector/pkg/types"
	"fmt"
)

func (c *Coordinator) Getter(name string) (types.GetterInfo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	state, ok := c.getters[name]

	if !ok {
		return types.GetterInfo{}, fmt.Errorf("such getter doesn't exist")
	} else {
		return types.GetterInfo{Name: name, State: state.State().String()}, nil
	}

}

func (c *Coordinator) GetterList() []types.GetterInfo {

	c.mu.Lock()
	defer c.mu.Unlock()

	res := make([]types.GetterInfo, 0)
	for name, getter := range c.getters {
		res = append(res, types.GetterInfo{Name: name, State: getter.State().String()})
	}

	return res

}

func (c *Coordinator) Start(name string) error {

	c.mu.Lock()
	_, ok := c.getters[name]
	c.mu.Unlock()
	if ok {
		return fmt.Errorf("already exist")
	}

	err := c.addGetter(name)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.getters[name].SetState(getters.Running)
	go c.getters[name].Start()
	c.mu.Unlock()

	return nil
}

func (c *Coordinator) StartAll() error {
	for _, name := range c.getterNames {
		c.mu.Lock()
		_, ok := c.getters[name]
		c.mu.Unlock()
		if ok {
			return fmt.Errorf("already exist")
		}

		err := c.addGetter(name)
		if err != nil {
			return err
		}

		c.mu.Lock()
		c.getters[name].SetState(getters.Running)
		go c.getters[name].Start()
		c.mu.Unlock()
	}
	return nil
}

func (c *Coordinator) StopAll() {
	for _, getter := range c.getters {
		getter.SetState(getters.Down)
	}
}

func (c *Coordinator) Stop(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	getter, ok := c.getters[name]
	if !ok {
		return fmt.Errorf("such getter don't exist")
	}
	getter.SetState(getters.Down)
	return nil
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
