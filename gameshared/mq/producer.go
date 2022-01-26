package mq

import "errors"

type Producer interface {
	PublishAsync(msg *Msg)
	Close()
}

func NewProducer(cfg Config, options ...Option) (Producer, error) {
	if !cfg.Open {
		return newNullProducer(cfg.Address)
	}

	opts := GetDefaultOptions()
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, err
		}
	}

	switch cfg.Type {
	case "kafka":
		return newKafkaProducer(cfg.Address, opts)
	case "nats":
		return nil, errors.New("not support type")
	case "nats_js":
		return nil, errors.New("not support type")
	default:
		return nil, errors.New("not support type")
	}
}
