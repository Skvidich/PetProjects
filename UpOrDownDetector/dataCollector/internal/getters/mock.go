package getters

import (
	"dataCollector/pkg/types"
	"errors"
	"math/rand"
	"time"
)

func Mock1() (types.StatusResponse, error) {

	if rand.Intn(10) == 0 {
		return types.StatusResponse{}, errors.New("service unavailable")
	}

	services := []string{"API Gateway", "Database", "Auth Service", "Cache", "Payment Processor"}
	statuses := []string{"OK", "Error", "Degraded", "Maintenance"}
	components := []string{"Storage", "Network", "CPU", "Memory", "API"}

	// Generate random components
	var randomComponents []types.Component
	for i := 0; i < rand.Intn(3)+1; i++ { // 1-3 components
		randomComponents = append(randomComponents, types.Component{
			Name:   components[rand.Intn(len(components))],
			Status: statuses[rand.Intn(len(statuses))],
		})
	}

	return types.StatusResponse{
		Name:       services[rand.Intn(len(services))],
		Time:       time.Now().UTC(),
		Components: randomComponents,
	}, nil
}

func Mock2() (types.StatusResponse, error) {

	if rand.Intn(5) == 0 {
		return types.StatusResponse{}, errors.New("connection timeout")
	}

	services := []string{"User Service", "Analytics", "Notification System", "File Storage", "AI Engine"}
	statuses := []string{"Active", "Inactive", "Overloaded", "Recovering", "Unknown"}
	components := []string{"Logs", "Security", "Load Balancer", "Backup", "Monitoring"}

	var randomComponents []types.Component
	for i := 0; i < rand.Intn(3)+2; i++ {
		randomComponents = append(randomComponents, types.Component{
			Name:   components[rand.Intn(len(components))],
			Status: statuses[rand.Intn(len(statuses))],
		})
	}

	return types.StatusResponse{
		Name:       services[rand.Intn(len(services))],
		Time:       time.Now().UTC(),
		Components: randomComponents,
	}, nil
}
