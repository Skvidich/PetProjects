package StatusRelay

import (
	"context"
	"dataCollector/common"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type StatusProducer struct {
	Ctx       context.Context
	CtxCancel context.CancelFunc
	Producer  sarama.AsyncProducer
	Topic     string
}

func NewStatusProducer(conf *sarama.Config, addr []string, topic string) (*StatusProducer, error) {

	producer, err := sarama.NewAsyncProducer(addr, conf)
	if err != nil {
		return nil, fmt.Errorf("can't create kafka producer %v", err)
	}
	ctx, ctxCancel := context.WithCancel(context.Background())
	return &StatusProducer{
		Ctx:       ctx,
		CtxCancel: ctxCancel,
		Producer:  producer,
		Topic:     topic,
	}, nil
}

func (prod *StatusProducer) Produce(mess common.StatusResponse) error {

	rawVal, err := json.Marshal(mess)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: prod.Topic,
		Key:   sarama.StringEncoder(mess.Name),
		Value: sarama.ByteEncoder(rawVal),
	}

	select {
	case prod.Producer.Input() <- msg:
		return nil
	case <-prod.Ctx.Done():
		return fmt.Errorf("producer context cancelled")
	}

}

func (prod *StatusProducer) Shutdown() error {
	prod.CtxCancel()
	err := prod.Producer.Close()
	if err != nil {
		return fmt.Errorf("error at closing producer: %v", err)
	}
	return nil
}
