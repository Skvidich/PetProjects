package StatusGetters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var GetterFuncList = map[string]StatusGetterFunc{
	"Github": GithubGet,
}

func GithubGet() (StatusResponse, error) {
	resp, err := http.Get("https://www.githubstatus.com/api/v2/summary.json")

	if err != nil {
		return StatusResponse{}, fmt.Errorf("error while get %v", err)
	}

	var data struct {
		Components []Component
	}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&data)
	if err != nil {
		return StatusResponse{}, fmt.Errorf("error decoding json %v", err)
	}

	res := StatusResponse{Name: "GitHub", Components: data.Components, Time: time.Now()}
	return res, nil

}
