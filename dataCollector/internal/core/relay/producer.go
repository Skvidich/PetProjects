package relay

import (
	"context"
	"dataCollector/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type Producer struct {
	ctx    context.Context
	cancel context.CancelFunc
	prod   sarama.AsyncProducer
	topic  string
}

func NewProducer(cfg *sarama.Config, brokers []string, topic string) (*Producer, error) {

	p, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create kafka producer %v", err)
	}
	ctx, ctxCancel := context.WithCancel(context.Background())
	return &Producer{
		ctx:    ctx,
		cancel: ctxCancel,
		prod:   p,
		topic:  topic,
	}, nil
}

func (p *Producer) Produce(mess types.ServiceStatus) error {

	raw, err := json.Marshal(mess)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(mess.Name),
		Value: sarama.ByteEncoder(raw),
	}

	select {
	case p.prod.Input() <- msg:
		return nil
	case <-p.ctx.Done():
		return context.Canceled
	}

}

func (p *Producer) Close() error {
	p.cancel()
	err := p.prod.Close()
	if err != nil {
		return fmt.Errorf("error at closing producer: %v", err)
	}
	return nil
}
