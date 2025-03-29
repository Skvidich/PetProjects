package statusConsumer

import (
	"context"
	"dataProcessor/common"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"sync"
)

const BUFFSIZE = 10

type Consumer interface {
	GetStatusMessage() (common.StatusMessage, error)
}

type StatusConsumer struct {
	Consumer       sarama.Consumer
	Ctx            context.Context
	CtxCancel      context.CancelFunc
	PartitionsWait sync.WaitGroup
	Output         chan common.StatusMessage
}

func NewStatusConsumer(config *sarama.Config, addrs []string) (*StatusConsumer, error) {
	consumer, err := sarama.NewConsumer(addrs, config)
	if err != nil {
		return nil, err
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	return &StatusConsumer{
		Consumer:       consumer,
		Ctx:            ctx,
		CtxCancel:      ctxCancel,
		Output:         make(chan common.StatusMessage, BUFFSIZE),
		PartitionsWait: sync.WaitGroup{},
	}, nil

}

func (cons *StatusConsumer) Run(topic string) {

	partitions, err := cons.Consumer.Partitions(topic)
	if err != nil {
		return
	}
	for _, partition := range partitions {
		var partCons sarama.PartitionConsumer
		partCons, err = cons.Consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			continue
		}
		cons.PartitionsWait.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer func() {
				err := pc.Close()
				if err != nil {
					log.Print(err.Error())
				}
				cons.PartitionsWait.Done()
			}()

			var statusMess common.StatusMessage
			for {
				select {
				case msg, ok := <-pc.Messages():
					if !ok {
						return
					}
					fmt.Printf("Партиция: %d, Смещение: %d, Ключ: %s, Значение: %s\n",
						msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
					err := json.Unmarshal(msg.Value, &statusMess)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					cons.Output <- statusMess
				case <-cons.Ctx.Done():
					return
				}

			}

		}(partCons)
	}

}

func (cons *StatusConsumer) Close() error {

	cons.CtxCancel()
	cons.PartitionsWait.Wait()
	err := cons.Consumer.Close()
	if err != nil {
		return err
	}
	close(cons.Output)
	return nil
}
func (cons *StatusConsumer) GetStatusMessage() (common.StatusMessage, error) {
	return common.StatusMessage{}, nil
}
