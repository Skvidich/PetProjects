package databaseSaver

import (
	"dataProcessor/common"
	"fmt"
)

type Saver interface {
	SaveStatistic(overview *common.StatisticOverview, components []common.ComponentStatistic) error
	SaveIncident(stat *common.Incident) error
}

type DbController struct {
}

func (cntrl *DbController) SaveStatistic(overview *common.StatisticOverview, components []common.ComponentStatistic) error {
	fmt.Println("component saved")
	return nil
}

func (cntrl *DbController) SaveIncident(stat *common.Incident) error {
	fmt.Println("incident saved")
	return nil
}
