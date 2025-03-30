package reader

import (
	"context"
	"dataProcessor/common"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"sync"
)

type KafkaConsumer struct {
	client       sarama.Consumer
	ctx          context.Context
	cancel       context.CancelFunc
	partitionsWg sync.WaitGroup
	queue        *ConcurrentQueue
	errMx        sync.Mutex
	lastErr      error
}

func NewConsumer(cfg *sarama.Config, brokers []string, queue *ConcurrentQueue) (*KafkaConsumer, error) {
	client, err := sarama.NewConsumer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	if queue == nil {
		queue = NewConcurrentQueue()
	}
	return &KafkaConsumer{
		client:       client,
		ctx:          ctx,
		cancel:       cancel,
		partitionsWg: sync.WaitGroup{},
		queue:        queue,
	}, nil

}

func (c *KafkaConsumer) ConsumeTopic(topic string) {

	partitions, err := c.client.Partitions(topic)
	if err != nil {
		c.handleError(fmt.Errorf("get partitions failed: %w", err))
		return
	}

	for _, partition := range partitions {
		var partitionConsumer sarama.PartitionConsumer
		partitionConsumer, err = c.client.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			c.handleError(err)
			continue
		}
		c.partitionsWg.Add(1)
		go c.handlePartition(partitionConsumer)
	}

}

func (c *KafkaConsumer) handlePartition(pc sarama.PartitionConsumer) {
	defer func() {
		err := pc.Close()
		if err != nil {
			c.handleError(err)
		}
		c.partitionsWg.Done()
	}()

	for {
		select {
		case msg, ok := <-pc.Messages():
			if !ok {
				return
			}
			var status common.ServiceStatus
			if err := json.Unmarshal(msg.Value, &status); err != nil {
				c.handleError(fmt.Errorf("message decode error: %w", err))
				continue
			}
			c.queue.Enqueue(&status)

		case <-c.ctx.Done():
			return
		}
	}
}

func (c *KafkaConsumer) handleError(err error) {
	log.Print(err.Error())
	// Some error checking

	c.errMx.Lock()
	defer c.errMx.Unlock()
	if c.lastErr != nil {
		c.lastErr = err
	}

}

func (c *KafkaConsumer) Close() error {

	c.cancel()
	c.partitionsWg.Wait()
	err := c.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *KafkaConsumer) BackupQueue() *ConcurrentQueue {
	return c.queue
}

func (c *KafkaConsumer) Next() (*common.ServiceStatus, error) {

	c.errMx.Lock()
	err := c.lastErr
	c.errMx.Unlock()

	if err != nil {
		return nil, err
	}
	mess := c.queue.Dequeue()

	return mess, nil
}
