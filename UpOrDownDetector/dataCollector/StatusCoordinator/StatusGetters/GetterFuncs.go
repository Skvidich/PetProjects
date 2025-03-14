package StatusGetters

import (
	"dataCollector/common"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var GetterFuncList = map[string]StatusGetterFunc{
	"Github":     GetterStatusWrapper("Github", "https://www.githubstatus.com/api/v2/summary.json"),
	"DropBox":    GetterStatusWrapper("DropBox", "https://status.dropbox.com/api/v2/summary.json"),
	"Discord":    GetterStatusWrapper("Discord", "https://status.discord.com/api/v2/summary.json"),
	"Cloudflare": GetterStatusWrapper("Cloudflare", "https://www.cloudflarestatus.com/api/v2/summary.json"),
	"Mock1":      Mock1Get,
	"Mock2":      Mock2Get,
}

func GetterStatusWrapper(name string, url string) func() (common.StatusResponse, error) {
	return func() (common.StatusResponse, error) {
		resp, err := http.Get(url)

		if err != nil {
			return common.StatusResponse{}, fmt.Errorf("error while get %v", err)
		}

		var data struct {
			Components []common.Component
		}
		decoder := json.NewDecoder(resp.Body)
		decoder.UseNumber()
		err = decoder.Decode(&data)
		if err != nil {
			return common.StatusResponse{}, fmt.Errorf("error decoding json %v", err)
		}

		res := common.StatusResponse{Name: name, Components: data.Components, Time: time.Now()}
		return res, nil
	}
}

func Mock1Get() (common.StatusResponse, error) {

	if rand.Intn(10) == 0 {
		return common.StatusResponse{}, errors.New("service unavailable")
	}

	services := []string{"API Gateway", "Database", "Auth Service", "Cache", "Payment Processor"}
	statuses := []string{"OK", "Error", "Degraded", "Maintenance"}
	components := []string{"Storage", "Network", "CPU", "Memory", "API"}

	// Generate random components
	var randomComponents []common.Component
	for i := 0; i < rand.Intn(3)+1; i++ { // 1-3 components
		randomComponents = append(randomComponents, common.Component{
			Name:   components[rand.Intn(len(components))],
			Status: statuses[rand.Intn(len(statuses))],
		})
	}

	return common.StatusResponse{
		Name:       services[rand.Intn(len(services))],
		Time:       time.Now().UTC(),
		Components: randomComponents,
	}, nil
}

func Mock2Get() (common.StatusResponse, error) {

	if rand.Intn(5) == 0 {
		return common.StatusResponse{}, errors.New("connection timeout")
	}

	services := []string{"User Service", "Analytics", "Notification System", "File Storage", "AI Engine"}
	statuses := []string{"Active", "Inactive", "Overloaded", "Recovering", "Unknown"}
	components := []string{"Logs", "Security", "Load Balancer", "Backup", "Monitoring"}

	var randomComponents []common.Component
	for i := 0; i < rand.Intn(3)+2; i++ {
		randomComponents = append(randomComponents, common.Component{
			Name:   components[rand.Intn(len(components))],
			Status: statuses[rand.Intn(len(statuses))],
		})
	}

	return common.StatusResponse{
		Name:       services[rand.Intn(len(services))],
		Time:       time.Now().UTC(),
		Components: randomComponents,
	}, nil
}
