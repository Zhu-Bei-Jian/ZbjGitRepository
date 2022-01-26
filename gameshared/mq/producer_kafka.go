package mq

import (
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type kafkaProducer struct {
	sarama.AsyncProducer
	wg    sync.WaitGroup
	close bool

	opts Options
}

func newKafkaProducer(address []string, opts Options) (*kafkaProducer, error) {

	logrus.WithFields(logrus.Fields{
		"url": address,
	}).Info("正在连接 kafka")

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Timeout = 5 * time.Second
	cfg.Net.DialTimeout = 5 * time.Second
	cfg.Net.MaxOpenRequests = 5

	p := &kafkaProducer{
		opts: opts,
	}

	var err error
	p.AsyncProducer, err = sarama.NewAsyncProducer(address, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": address,
		}).WithError(err).Error("connect kafka")
		return nil, err
	}

	p.Run()
	return p, nil
}

// Run 运行.
func (p *kafkaProducer) Run() {
	p.errHandler()
}

// Close 关闭.
func (p *kafkaProducer) Close() {
	p.close = true
	p.AsyncProducer.Close()
	p.wg.Wait()
}

func (p *kafkaProducer) errHandler() {
	p.wg.Add(1)
	SafeGo(func() {
		defer p.wg.Done()

		err := p.Errors()
		success := p.Successes()
		for {
			select {
			case err, ok := <-err:
				if !ok {
					return
				}

				var data []byte
				if v, ok := err.Msg.Value.(sarama.ByteEncoder); ok {
					data = v
				}
				topic := err.Msg.Topic
				if p.opts.aecb != nil {
					p.opts.aecb(&Msg{
						Topic: topic,
						Data:  data,
					}, err)
				} else {
					logrus.WithFields(logrus.Fields{
						"topic": topic,
						"data":  data,
					}).WithError(err).Error("kafkaProducer pub")
				}
			case _, ok := <-success:
				if !ok {
					return
				}

				//str := ""
				//value, ok := msg.Value.(sarama.ByteEncoder)
				//if ok {
				//	str = string(value)
				//}
				//logrus.WithFields(logrus.Fields{
				//	"topic":     msg.Topic,
				//	"msg":       str,
				//	"offset":    msg.Offset,
				//	"partition": msg.Partition,
				//}).Debugf("kafkaProducer pub message success")
			}
		}
	})
}

func (p *kafkaProducer) PublishAsync(m *Msg) {
	msg := &sarama.ProducerMessage{
		Topic: m.Topic,
		Value: sarama.ByteEncoder(m.Data),
		Key:   sarama.ByteEncoder(""),
	}
	p.Input() <- msg
}
