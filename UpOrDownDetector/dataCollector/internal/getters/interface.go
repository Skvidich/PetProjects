package getters

import (
	"dataCollector/pkg/types"
	"sync"
	"time"
)

type State byte

const (
	Running State = iota
	Down
	Error
)

func (st State) String() string {
	var res string
	switch st {
	case Running:
		res = "Running"
	case Down:
		res = "Down"
	case Error:
		res = "Error"
	}
	return res
}

const outBuffSize = 2

type GetterFunc func() (types.ServiceStatus, error)

type Feedback struct {
	Err  error
	Name string
}

type Getter struct {
	name         string
	interval     time.Duration
	state        State
	mu           sync.Mutex
	OutChan      chan types.ServiceStatus
	feedbackChan *chan Feedback
	get          GetterFunc
}

func New(name string, get GetterFunc, interval time.Duration, feedbackChan *chan Feedback) *Getter {
	return &Getter{
		name:         name,
		interval:     interval,
		state:        Down,
		mu:           sync.Mutex{},
		OutChan:      make(chan types.ServiceStatus, outBuffSize),
		feedbackChan: feedbackChan,
		get:          get,
	}
}

func (g *Getter) State() State {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.state
}

func (g *Getter) SetState(st State) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.state = st
}

func (g *Getter) Start() {

	var err error
	var resp types.ServiceStatus
	defer close(g.OutChan)
	for {

		if g.State() == Down {
			*g.feedbackChan <- Feedback{Name: g.name, Err: nil}
			return
		}
		resp, err = g.get()

		if err != nil {
			g.SetState(Error)
			*g.feedbackChan <- Feedback{Name: g.name, Err: err}
			return
		}

		g.OutChan <- resp

		time.Sleep(g.interval)
	}
}

var Getters = map[string]GetterFunc{
	"Github":     GetterWrap("Github", "https://www.githubstatus.com/api/v2/summary.json"),
	"DropBox":    GetterWrap("DropBox", "https://status.dropbox.com/api/v2/summary.json"),
	"Discord":    GetterWrap("Discord", "https://status.discord.com/api/v2/summary.json"),
	"Cloudflare": GetterWrap("Cloudflare", "https://www.cloudflarestatus.com/api/v2/summary.json"),
	"Mock1":      Mock1,
	"Mock2":      Mock2,
}
