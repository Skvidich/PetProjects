package coordinator

import (
	"dataCollector/internal/core/getters"
	"dataCollector/internal/logger"
	"dataCollector/pkg/types"
	"fmt"
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
	errLog       logger.Logger
}

func New(interval time.Duration, names []string, errLog logger.Logger) *Coordinator {

	out := make(chan types.ServiceStatus)

	c := &Coordinator{
		getterNames:  names,
		mu:           sync.Mutex{},
		getters:      make(map[string]*getters.Getter, len(names)),
		OutChan:      out,
		feedbackChan: make(chan getters.Feedback, len(names)),
		interval:     interval,
		errLog:       errLog,
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
			c.errLog.LogError(feedback.Name+" got error", feedback.Err)
		} else {
			c.removeGetter(feedback.Name)
		}

	}
}

func (c *Coordinator) removeGetter(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.getters, name)

}

func (c *Coordinator) addGetter(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	getFunc, ok := getters.Getters[name]
	if !ok {
		return fmt.Errorf("such getter don't exist")
	}

	c.getters[name] = getters.New(name, getFunc, c.interval, &c.feedbackChan)
	go c.forward(c.getters[name].OutChan)
	return nil
}

func (c *Coordinator) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.getters)

}
