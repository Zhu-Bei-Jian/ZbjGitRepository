package gameutil

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

// NatsClient ...
type NatsClient struct {
	conn *nats.Conn
	subs map[string]*nats.Subscription
}

func (p *NatsClient) Conn() *nats.Conn {
	return p.conn
}

// Init ...
func (p *NatsClient) Init(address string) error {
	var err error
	url := nats.DefaultURL
	if address != "" {
		url = address
	}
	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Info("正在连接 nats")

	p.conn, err = nats.Connect(
		url,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		p.onDisconnect(),
		p.onReconnect(),
		p.onClosed(),
	)

	if err != nil {
		return err
	}

	p.subs = make(map[string]*nats.Subscription)

	logrus.Info("连接 nats 成功")

	return nil
}

// Close ...
func (p *NatsClient) Close() error {
	for _, v := range p.subs {
		err := v.Unsubscribe()
		if err != nil {
			return err
		}
	}
	p.conn.Close()
	return nil
}

// Pub ...
func (p *NatsClient) Pub(topic string, msg []byte) error {
	err := p.conn.Publish(topic, msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic": topic,
		}).WithError(err).Error("NatsClient Pub error")
	}
	return err
}

// Sub ...
func (p *NatsClient) Sub(topic string, cb func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, ok := p.subs[topic]
	if ok {
		return sub, nil
	}
	sub, err := p.conn.Subscribe(topic, cb)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic": topic,
		}).WithError(err).Error("NatsClient Sub error")
		return nil, err
	}
	return sub, nil
}

//chan Sub
func (p *NatsClient) ChanSub(topic string, ch chan *nats.Msg) (*nats.Subscription, error) {
	sub, ok := p.subs[topic]
	if ok {
		return sub, nil
	}
	sub, err := p.conn.ChanSubscribe(topic, ch)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic": topic,
		}).WithError(err).Error("NatsClient Sub error")
		return nil, err
	}
	return sub, nil
}

// Unsubscribe ...
func (p *NatsClient) Unsubscribe(topic string) error {
	sub, ok := p.subs[topic]
	if ok {
		return sub.Unsubscribe()
	}
	return nil
}

func (p *NatsClient) onDisconnect() nats.Option {
	return nats.DisconnectHandler(func(*nats.Conn) {
		logrus.Warn("nats server disconnected")
	})
}

func (p *NatsClient) onReconnect() nats.Option {
	return nats.ReconnectHandler(func(*nats.Conn) {
		logrus.Warn("nats server reconnected")
	})
}

func (p *NatsClient) onClosed() nats.Option {
	return nats.ClosedHandler(func(*nats.Conn) {
		logrus.Warn("nats server closed")
	})
}
