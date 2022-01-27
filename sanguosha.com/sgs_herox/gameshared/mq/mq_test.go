package mq

import (
	"fmt"
	"testing"
	"time"
)

func TestNewProducer_Null(t *testing.T) {
	c := Config{
		Open: false,
		Type: "",
	}

	p, err := NewProducer(c, PublishAsyncErrHandler(func(msg *Msg, err error) {
		if err != nil {
			t.Fatal(err)
		}
	}))

	if err != nil {
		t.Fatal(err)
	}

	p.PublishAsync(nil)
}

func TestMQ_Kafka(t *testing.T) {
	cfg := Config{
		Open:    true,
		Type:    "kafka",
		Address: []string{"10.225.254.248:9092"},
	}

	topic := "test"
	data := "hello"

	c, err := NewConsumer(cfg)
	if err != nil {
		t.Fatal(err)
	}
	c.Sub(topic, func(v []byte) {
		fmt.Println(topic, string(v))
		if string(v) != data {
			t.Fatalf("expect:%s get:%s", data, string(v))
		}
	})

	p, err := NewProducer(cfg, PublishAsyncErrHandler(func(msg *Msg, err error) {
		if err != nil {
			t.Fatal(err)
		}
	}))

	if err != nil {
		t.Fatal(err)
	}

	p.PublishAsync(&Msg{
		Topic: topic,
		Data:  []byte(data),
	})

	p.Close()

	time.Sleep(time.Second * 60)
}
