package getters

import (
	"dataCollector/pkg/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetterWrap(name string, url string) func() (types.ServiceStatus, error) {
	return func() (types.ServiceStatus, error) {
		resp, err := http.Get(url)

		if err != nil {
			return types.ServiceStatus{}, fmt.Errorf("error while get %v", err)
		}
		defer resp.Body.Close()
		var data struct {
			Components []types.Component
		}

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return types.ServiceStatus{}, fmt.Errorf("error decoding json %v", err)
		}

		return types.ServiceStatus{
			Name:       name,
			Components: data.Components,
			Time:       time.Now(),
		}, nil
	}
}
