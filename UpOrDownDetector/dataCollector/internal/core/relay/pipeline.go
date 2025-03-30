package relay

import (
	"dataCollector/internal/utils"
	"dataCollector/pkg/types"
	"github.com/IBM/sarama"
)

type Relay struct {
	inChan    chan types.ServiceStatus
	logging   bool
	resend    bool
	done      chan struct{}
	processor StatusProcess
	producer  *Producer
}

func NewRelay(in chan types.ServiceStatus, logging bool, resend bool) *Relay {

	processStatus := func(response types.ServiceStatus) error { return nil }

	return &Relay{
		inChan:    in,
		logging:   logging,
		resend:    resend,
		done:      make(chan struct{}),
		processor: processStatus,
	}
}

func (r *Relay) SetupProducer(topic string, brokers []string) error {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3
	cfg.Net.MaxOpenRequests = 5
	cfg.Producer.Partitioner = sarama.NewHashPartitioner

	var err error
	r.producer, err = NewProducer(cfg, brokers, topic)
	if err != nil {
		return err
	}
	return nil
}

func (r *Relay) InitPipeline() {
	if r.logging {
		r.wrapLogger()
	}
	if r.resend {
		r.wrapResend()
	}
}

func (r *Relay) process() {
	for status := range r.inChan {

		err := r.processor(status)
		if err != nil {
			break
		}
	}
	r.done <- struct{}{}
}

type StatusProcess func(st types.ServiceStatus) error

func (r *Relay) wrapLogger() {
	originalProcessStatus := r.processor
	r.processor = func(status types.ServiceStatus) error {
		utils.LogStatus(status)
		err := originalProcessStatus(status)
		if err != nil {
			return err
		}
		return nil
	}
}

func (r *Relay) wrapResend() {
	orig := r.processor

	r.processor = func(status types.ServiceStatus) error {
		err := r.producer.Produce(status)
		if err != nil {
			return err
		}

		return orig(status)
	}
}
