package alert

import (
	"context"
	"dataProcessor/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"time"
)

type KafkaReporter struct {
	client     sarama.AsyncProducer
	ctx        context.Context
	cancel     context.CancelFunc
	startTopic string
	endTopic   string
}

func NewKafkaReporter(cfg *sarama.Config, brokers []string, startTopic string, endTopic string) (*KafkaReporter, error) {

	client, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create kafka producer %v", err)
	}
	ctx, ctxCancel := context.WithCancel(context.Background())
	return &KafkaReporter{
		ctx:        ctx,
		cancel:     ctxCancel,
		client:     client,
		startTopic: startTopic,
		endTopic:   endTopic,
	}, nil

}

type IncidentStart struct {
	Service   string    `json:"service"`
	StartTime time.Time `json:"start_time"`
	Component string    `json:"component"`
	Status    string    `json:"status"`
}

func (r *KafkaReporter) NotifyStart(service string, startTime time.Time, component models.Component) error {
	raw, err := json.Marshal(IncidentStart{
		Service:   service,
		StartTime: startTime,
		Component: component.Name,
		Status:    component.Status,
	})
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: r.startTopic,
		Key:   sarama.StringEncoder(service),
		Value: sarama.ByteEncoder(raw),
	}

	select {
	case r.client.Input() <- msg:
		return nil
	case <-r.ctx.Done():
		return context.Canceled
	}

}

func (r *KafkaReporter) NotifyEnd(incident *models.ServiceIncident) error {
	raw, err := json.Marshal(*incident)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: r.endTopic,
		Key:   sarama.StringEncoder(incident.Name),
		Value: sarama.ByteEncoder(raw),
	}

	select {
	case r.client.Input() <- msg:
		return nil
	case <-r.ctx.Done():
		return context.Canceled
	}
}

func (r *KafkaReporter) Close() error {
	r.cancel()
	err := r.client.Close()
	if err != nil {
		return fmt.Errorf("error at closing producer: %v", err)
	}
	return nil
}
