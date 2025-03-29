package databaseController

import (
	"dataProcessor/common"
	"fmt"
)

type StatisticSaver interface {
	SaveOverview(stat *common.StatisticOverview) (bool, error)
	SaveComponent(stat *common.ComponentStatistic) (bool, error)
}

type IncidentSaver interface {
	SaveIncident(stat *common.StatisticOverview) (bool, error)
}

type DbController struct {
}

func (cntrl *DbController) SaveOverview(stat *common.StatisticOverview) (bool, error) {
	fmt.Println("overview saved")
	return false, nil
}

func (cntrl *DbController) SaveComponent(stat *common.ComponentStatistic) (bool, error) {
	fmt.Println("component saved")
	return false, nil
}

func (cntrl *DbController) SaveIncident(stat *common.StatisticOverview) (bool, error) {
	fmt.Println("incident saved")
	return false, nil
}
