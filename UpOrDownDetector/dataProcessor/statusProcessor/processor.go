package statusProcessor

import (
	"context"
	"dataProcessor/common"
	"dataProcessor/databaseSaver"
	"dataProcessor/incidentNotificator"
	"dataProcessor/statusConsumer"
	"time"
)

const StatisticDeadline = time.Hour

type incidentKey struct {
	Name      string
	Component string
}
type incidentValue struct {
	Status    string
	StartTime time.Time
}

type aggregatedStat struct {
	Name       string
	StartTime  time.Time
	EndTime    time.Time
	Components map[string]map[string]int // component -> status -> count
}

type StatsBuffer map[string]*aggregatedStat
type IncidentBuffer map[incidentKey]incidentValue
type StatusProcessor struct {
	lastIncidents IncidentBuffer
	currStatistic StatsBuffer
	saver         databaseSaver.Saver
	consumer      statusConsumer.Consumer
	notificator   incidentNotificator.Notificator
	Ctx           context.Context
	CtxCancel     context.CancelFunc
	Quit          chan common.ErrorInfo
}

func NewStatusProcessor(
	saver databaseSaver.Saver,
	consumer statusConsumer.Consumer,
	notificator incidentNotificator.Notificator,
	parent context.Context,
) (*StatusProcessor, error) {

	ctx, cancel := context.WithCancel(parent)

	return &StatusProcessor{
		lastIncidents: make(map[incidentKey]incidentValue),
		currStatistic: make(map[string]*aggregatedStat),
		saver:         saver,
		consumer:      consumer,
		notificator:   notificator,
		Ctx:           ctx,
		CtxCancel:     cancel,
		Quit:          make(chan common.ErrorInfo),
	}, nil
}

func (proc *StatusProcessor) Run(timeout time.Duration) {

	ctx, _ := context.WithTimeout(context.Background(), timeout)

	for {
		select {
		case <-proc.Ctx.Done():
			return
		case <-ctx.Done():
			mess, err := proc.consumer.GetStatusMessage()
			if err != nil {
				proc.SendQuitMessage(common.ErrorInfo{Entity: common.ConsumerError, Err: err})
				return
			}
			errInf := proc.ProcessMessage(mess)
			if errInf.Err != nil {
				proc.SendQuitMessage(errInf)
				return
			}
			ctx, _ = context.WithTimeout(context.Background(), timeout)

		}
	}

}

func (proc *StatusProcessor) ProcessMessage(mess common.StatusMessage) common.ErrorInfo {
	var errInfo common.ErrorInfo

	errInfo = proc.checkIncidents(mess)
	if errInfo.Err != nil {
		return errInfo
	}

	errInfo = proc.saveStatistic(mess)
	if errInfo.Err != nil {
		return errInfo
	}

	return common.ErrorInfo{Err: nil}
}

func (proc *StatusProcessor) checkIncidents(mess common.StatusMessage) common.ErrorInfo {
	var err error

	for _, component := range mess.Components {
		key := incidentKey{
			Name:      mess.Name,
			Component: component.Name,
		}
		incident, ok := proc.lastIncidents[key]

		if ok && common.IsUpStatus(component.Status) {
			endedIncident := &common.Incident{
				Name:          key.Name,
				ComponentName: key.Name,
				Status:        incident.Status,
				StartTime:     incident.StartTime,
				EndTime:       mess.Time,
			}
			err = proc.saver.SaveIncident(endedIncident)
			if err != nil {
				return common.ErrorInfo{Entity: common.SaverError, Err: err}
			}

			err = proc.notificator.EndIncident(endedIncident)
			if err != nil {
				return common.ErrorInfo{Entity: common.NotificationError, Err: err}
			}
			delete(proc.lastIncidents, key)
		}

		if !ok && common.IsDownStatus(component.Status) {
			proc.lastIncidents[incidentKey{Name: mess.Name, Component: component.Name}] = incidentValue{StartTime: mess.Time, Status: component.Status}
			err = proc.notificator.StartIncident(mess.Name, mess.Time, component)
			if err != nil {
				return common.ErrorInfo{Entity: common.NotificationError, Err: err}
			}
		}
	}
	return common.ErrorInfo{Err: nil}
}

func (proc *StatusProcessor) saveStatistic(mess common.StatusMessage) common.ErrorInfo {

	_, ok := proc.currStatistic[mess.Name]
	if !ok {
		proc.currStatistic[mess.Name] = &aggregatedStat{Name: mess.Name, StartTime: mess.Time}
	}

	for _, component := range mess.Components {
		proc.currStatistic[mess.Name].Components[component.Name][component.Status]++
	}

	deadline := proc.currStatistic[mess.Name].StartTime.Add(StatisticDeadline)
	if mess.Time.After(deadline) {
		proc.currStatistic[mess.Name].EndTime = mess.Time
		overview, components := ParseAggregated(proc.currStatistic[mess.Name])
		err := proc.saver.SaveStatistic(overview, components)
		if err != nil {
			return common.ErrorInfo{Entity: common.SaverError, Err: err}
		}
		delete(proc.currStatistic, mess.Name)
	}

	return common.ErrorInfo{Err: nil}
}

func (proc *StatusProcessor) SendQuitMessage(errInf common.ErrorInfo) {
	proc.Quit <- errInf
}

func (proc *StatusProcessor) GetBackup() (IncidentBuffer, StatsBuffer) {
	return proc.lastIncidents, proc.currStatistic
}

func (proc *StatusProcessor) Close() {
	close(proc.Quit)
}

func ParseAggregated(aggr *aggregatedStat) (*common.StatisticOverview, []common.ComponentStatistic) {
	overview := &common.StatisticOverview{
		Name:      aggr.Name,
		StartTime: aggr.StartTime,
		EndTime:   aggr.EndTime,
	}

	compStat := make([]common.ComponentStatistic, 0)

	var temp common.ComponentStatistic
	for name, stat := range aggr.Components {

		for status, count := range stat {
			temp.ComponentName = name
			temp.Status = status
			temp.Count = count
			compStat = append(compStat, temp)
		}

	}
	return overview, compStat
}
