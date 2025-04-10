package getters

import (
	"dataCollector/pkg/types"
	"math/rand"
	"time"
)

func Mock1() (types.ServiceStatus, error) {

	components := make([]types.Component, 2)

	components[0] = types.Component{
		Name:   "Something",
		Status: "operational",
	}

	components[1] = types.Component{
		Name:   "Another",
		Status: "operational",
	}

	return types.ServiceStatus{
		Name:       "Mock1",
		Time:       time.Now().UTC(),
		Components: components,
	}, nil
}

func Mock2() (types.ServiceStatus, error) {

	components := make([]types.Component, 2)

	components[0] = types.Component{
		Name: "Something",
	}

	if rand.Int()%11 < 7 {
		components[0].Status = "operational"
	} else {
		components[0].Status = "down"
	}
	components[1] = types.Component{
		Name: "Another",
	}

	if rand.Int()%11 < 7 {
		components[1].Status = "operational"
	} else {
		components[1].Status = "down"
	}
	return types.ServiceStatus{
		Name:       "Mock2",
		Time:       time.Now().UTC(),
		Components: components,
	}, nil
}
