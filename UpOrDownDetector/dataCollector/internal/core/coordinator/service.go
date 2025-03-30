package coordinator

import (
	"dataCollector/internal/getters"
	"dataCollector/internal/utils"
	"dataCollector/pkg/types"
	"sync"
	"time"
)

type Coordinator struct {
	getters      map[string]*getters.Getter
	mu           sync.Mutex
	getterNames  []string
	OutChan      chan types.ServiceStatus
	feedbackChan chan getters.Feedback
	interval     time.Duration
}

func New(interval time.Duration, names []string) *Coordinator {

	out := make(chan types.ServiceStatus)

	c := &Coordinator{
		getterNames:  names,
		mu:           sync.Mutex{},
		getters:      make(map[string]*getters.Getter, len(names)),
		OutChan:      out,
		feedbackChan: make(chan getters.Feedback, len(names)),
		interval:     interval,
	}
	go c.handleFeedback()
	return c
}

func (c *Coordinator) forward(src <-chan types.ServiceStatus) { // Было attachToBus
	for s := range src {
		c.OutChan <- s
	}
}

func (c *Coordinator) handleFeedback() {

	for feedback := range c.feedbackChan {

		if feedback.Err != nil {
			logError(feedback)
		} else {
			c.removeGetter(feedback.Name)
		}

	}
}

func logError(getErr getters.Feedback) {
	utils.LogError(getErr.Name + " " + getErr.Err.Error())
}

func (c *Coordinator) removeGetter(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.getters, name)

}

func (c *Coordinator) addGetter(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.getters[name] = getters.New(name, getters.Getters[name], c.interval, &c.feedbackChan)
	go c.forward(c.getters[name].OutChan)
}

func (c *Coordinator) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.getters)

}
