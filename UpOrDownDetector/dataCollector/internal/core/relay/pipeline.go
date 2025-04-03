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
	processor []StatusProcess
	producer  *Producer
	store     storage.Storage
	errLog    logger.Logger
}

func NewRelay(in chan types.ServiceStatus, save bool, resend bool, errLog logger.Logger, store storage.Storage) *Relay {

	rel := &Relay{
		inChan:    in,
		done:      make(chan struct{}),
		processor: make([]StatusProcess, 0),
		store:     store,
		errLog:    errLog,
	}
	if save {
		rel.addSave()
	}
	if resend {
		rel.addResend()
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
		for _, proc := range r.processor {
			err := proc(status)
			if err != nil {
				r.errLog.LogError("can't process", err)
			}
		}
	}
	r.done <- struct{}{}
}

type StatusProcess func(st types.ServiceStatus) error

func (r *Relay) addSave() {
	proc := func(status types.ServiceStatus) error {
		err := r.store.StoreRawReport(&status)
		return err
	}
	r.processor = append(r.processor, proc)
}

func (r *Relay) addResend() {
	proc := func(status types.ServiceStatus) error {
		err := r.producer.Produce(status)
		return err
	}
	r.processor = append(r.processor, proc)
}
