package StatusGetters

import (
	"dataCollector/common"
	"sync"
	"time"
)

type GetterState byte

// TO-DO
// add isPaused getter state
const (
	IsRunning GetterState = iota
	IsDown
	IsError
)

func (st GetterState) String() string {
	var res string
	switch st {
	case IsRunning:
		res = "IsRunning"
	case IsDown:
		res = "IsDown"
	case IsError:
		res = "IsError"
	}
	return res
}

const outBuffSize = 2

type StatusGetterFunc func() (common.StatusResponse, error)

type GetterFeedback struct {
	Err  error
	Name string
}

type StatusGetter struct {
	Name         string
	delay        time.Duration
	state        GetterState
	muState      sync.Mutex
	OutputChan   chan common.StatusResponse
	feedbackChan *chan GetterFeedback
	getFunc      StatusGetterFunc
}

func NewStatusGetter(name string, getFunc StatusGetterFunc, delay time.Duration, feedbackChan *chan GetterFeedback) *StatusGetter {
	return &StatusGetter{
		Name:         name,
		delay:        delay,
		state:        IsDown,
		muState:      sync.Mutex{},
		OutputChan:   make(chan common.StatusResponse, outBuffSize),
		feedbackChan: feedbackChan,
		getFunc:      getFunc,
	}
}

func (g *StatusGetter) GetState() GetterState {
	g.muState.Lock()
	defer g.muState.Unlock()
	return g.state
}

func (g *StatusGetter) SetState(st GetterState) {
	g.muState.Lock()
	defer g.muState.Unlock()
	g.state = st
}

func (g *StatusGetter) RunProcess() {

	var err error
	var resp common.StatusResponse
	defer close(g.OutputChan)
	for {

		if g.GetState() == IsDown {

			*g.feedbackChan <- GetterFeedback{Name: g.Name, Err: nil}
			return
		}
		resp, err = g.getFunc()

		if err != nil {
			g.SetState(IsError)
			*g.feedbackChan <- GetterFeedback{Name: g.Name, Err: err}
			return
		}

		g.OutputChan <- resp

		time.Sleep(g.delay)
	}
}
