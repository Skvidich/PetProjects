package relay

import (
	"dataCollector/internal/core/storage"
	"dataCollector/internal/logger"
	"dataCollector/pkg/types"
	"github.com/IBM/sarama"
)

type Relay struct {
	inChan    chan types.ServiceStatus
	done      chan struct{}
	processor StatusProcess
	producer  *Producer
	store     storage.Storage
	errLog    logger.Logger
}

func NewRelay(in chan types.ServiceStatus, save bool, resend bool, errLog logger.Logger, store storage.Storage) *Relay {

	rel := &Relay{
		inChan:    in,
		done:      make(chan struct{}),
		processor: func(response types.ServiceStatus) error { return nil },
		store:     store,
		errLog:    errLog,
	}
	if save {
		rel.wrapSave()
	}
	if resend {
		rel.wrapResend()
	}

	return rel
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

}

func (r *Relay) process() {
	for status := range r.inChan {

		err := r.processor(status)
		if err != nil {
			r.errLog.LogError("can't process", err)
		}
	}
	r.done <- struct{}{}
}

type StatusProcess func(st types.ServiceStatus) error

func (r *Relay) wrapSave() {
	orig := r.processor
	r.processor = func(status types.ServiceStatus) error {
		err := r.store.StoreRawReport(&status)
		if err != nil {
			return err
		}

		return orig(status)
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
