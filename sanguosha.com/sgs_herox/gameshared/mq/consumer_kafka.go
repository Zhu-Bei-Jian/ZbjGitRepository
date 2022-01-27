package mq

import (
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type kafkaConsumer struct {
	sarama.Consumer
	wg        sync.WaitGroup
	consumers []sarama.PartitionConsumer
	opts      Options
}

func newKafkaConsumer(address []string, opts Options) (*kafkaConsumer, error) {
	logrus.WithFields(logrus.Fields{
		"url": address,
	}).Info("正在连接 kafka")

	var err error
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Net.DialTimeout = 5 * time.Second
	cfg.Net.MaxOpenRequests = 5

	p := &kafkaConsumer{
		opts: opts,
	}
	p.Consumer, err = sarama.NewConsumer(address, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": address,
		}).WithError(err).Error("连接 kafka Consumer")
		return nil, err
	}
	return p, nil
}

// Close 关闭.
func (p *kafkaConsumer) Close() {
	for _, v := range p.consumers {
		v.Close()
	}
	p.Consumer.Close()
	p.wg.Wait()
}

// Sub 注册topic回调.
func (p *kafkaConsumer) Sub(topic string, callback func([]byte)) error {
	partitionList, err := p.Consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for _, v := range partitionList {
		consumer, err := p.Consumer.ConsumePartition(topic, v, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		p.consumers = append(p.consumers, consumer)
		p.wg.Add(1)

		SafeGo(func() {
			defer p.wg.Done()
			for {
				quit := false
				select {
				//接收消息通道和错误通道的内容.
				case msg := <-consumer.Messages():
					if msg == nil {
						quit = true
						break
					}

					logrus.WithFields(logrus.Fields{
						"offset": msg.Offset,
						"msg":    string(msg.Value),
						"topic":  msg.Topic,
					}).Debug("kafkaConsumer. sub")
					callback(msg.Value)
				case err := <-consumer.Errors():
					if err == nil {
						quit = true
						break
					}
					logrus.WithFields(logrus.Fields{
						"topic":     err.Topic,
						"partition": err.Partition,
					}).WithError(err).Error("kafkaConsumer Err")
				}
				if quit {
					break
				}
			}
		})
	}

	return nil
}
