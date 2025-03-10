package StatusGetters

import (
	"sync"
	"time"
)

type StatusResponse struct {
	Name       string      `json:"name"`
	Components []Component `json:"components"`
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type GetterState byte

const (
	IsRunning GetterState = iota

	IsDown
	IsError
)
const OutBuffSize = 2

type StatusGetterFunc func() (StatusResponse, error)

type GetterError struct {
	Err  error
	Name string
}

type StatusGetter struct {
	Name       string
	Delay      time.Duration
	State      GetterState
	MuState    sync.Mutex
	OutputChan chan StatusResponse
	ErrorChan  chan GetterError
	GetFunc    StatusGetterFunc
}

func NewStatusGetter(name string, getFunc StatusGetterFunc, delay time.Duration, errChan *chan GetterError) StatusGetter {
	return StatusGetter{
		Name:       name,
		Delay:      delay,
		State:      IsDown,
		MuState:    sync.Mutex{},
		OutputChan: make(chan StatusResponse, OutBuffSize),
		ErrorChan:  *errChan,
		GetFunc:    getFunc,
	}
}

func (g *StatusGetter) GetState() GetterState {
	g.MuState.Lock()
	defer g.MuState.Unlock()
	return g.State
}

func (g *StatusGetter) SetState(st GetterState) {
	g.MuState.Lock()
	defer g.MuState.Unlock()
	g.State = st
}

func (g *StatusGetter) RunProcess() {

	var err error
	var resp StatusResponse
	for {

		if g.GetState() == IsDown {
			return
		}
		resp, err = g.GetFunc()

		if err != nil {
			g.SetState(IsError)
			g.ErrorChan <- GetterError{Name: g.Name, Err: err}
			return
		}

		g.OutputChan <- resp

		time.Sleep(g.Delay)
	}
}
