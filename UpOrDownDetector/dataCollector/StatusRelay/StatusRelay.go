package StatusRelay

import (
	"dataCollector/common"
	"github.com/IBM/sarama"
)

type StatusRelay struct {
	inputChan     chan common.StatusResponse
	isLog         bool
	isResend      bool
	endChan       chan struct{}
	processStatus StatusProcess
	Producer      *StatusProducer
}

func NewStatusRelay(input chan common.StatusResponse, isLog bool, isResend bool) *StatusRelay {

	processStatus := func(response common.StatusResponse) error { return nil }

	return &StatusRelay{
		inputChan:     input,
		isLog:         isLog,
		isResend:      isResend,
		endChan:       make(chan struct{}),
		processStatus: processStatus,
	}
}
func (rel *StatusRelay) InitProducer(topic string, addrs []string) error {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 3
	conf.Net.MaxOpenRequests = 5
	conf.Producer.Partitioner = sarama.NewHashPartitioner

	var err error
	rel.Producer, err = NewStatusProducer(conf, addrs, topic)
	if err != nil {
		return err
	}
	return nil
}

func (rel *StatusRelay) InitProcess() {
	if rel.isLog {
		rel.logWrap()
	}
	if rel.isResend {
		rel.resendWrap()
	}
}

func (rel *StatusRelay) inputProcess() {
	for status := range rel.inputChan {

		err := rel.processStatus(status)
		if err != nil {
			break
		}
	}
	rel.endChan <- struct{}{}
}

type StatusProcess func(st common.StatusResponse) error

func (rel *StatusRelay) logWrap() {
	originalProcessStatus := rel.processStatus
	rel.processStatus = func(status common.StatusResponse) error {
		common.LogStatus(status)
		err := originalProcessStatus(status)
		if err != nil {
			return err
		}
		return nil
	}
}

func (rel *StatusRelay) resendWrap() {
	originalProcessStatus := rel.processStatus

	rel.processStatus = func(status common.StatusResponse) error {
		err := rel.Producer.Produce(status)
		if err != nil {
			return err
		}
		err = originalProcessStatus(status)
		if err != nil {
			return err
		}
		return nil
	}
}
