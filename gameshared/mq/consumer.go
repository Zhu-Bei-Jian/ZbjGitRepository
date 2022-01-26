package mq

import "errors"

type Consumer interface {
	Sub(topic string, cb func([]byte)) error
	Close()
}

func NewConsumer(cfg Config, options ...Option) (Consumer, error) {
	if !cfg.Open {
		return newNullConsumer(cfg.Address)
	}

	opts := GetDefaultOptions()
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, err
		}
	}

	switch cfg.Type {
	case "kafka":
		return newKafkaConsumer(cfg.Address, opts)
	case "nats":
		return nil, errors.New("not support type")
	case "nats_js":
		return nil, errors.New("not support type")
	default:
		return nil, errors.New("not support type")
	}
}
