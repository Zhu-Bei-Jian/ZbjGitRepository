package pio

import (
	"errors"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type poolProtoIO struct {
	pio  ProtoIO
	pool *Pool
	err  error
}

func (p *poolProtoIO) Read() (proto.Message, error) {
	msg, err := p.pio.Read()
	p.err = err
	return msg, err
}
func (p *poolProtoIO) ReadTimeout(d time.Duration) (proto.Message, error) {
	msg, err := p.pio.ReadTimeout(d)
	p.err = err
	return msg, err
}
func (p *poolProtoIO) Write(msg proto.Message) error {
	p.err = p.pio.Write(msg)
	return p.err
}
func (p *poolProtoIO) WriteTimeout(msg proto.Message, d time.Duration) error {
	p.err = p.pio.WriteTimeout(msg, d)
	return p.err
}
func (p *poolProtoIO) Close() error {
	if p.err != nil {
		p.err = p.pio.Close()
		return p.err
	}
	p.pool.put(p)
	return nil
}

type errorProtoIO struct {
	error
}

func (p errorProtoIO) Read() (proto.Message, error) {
	return nil, p
}
func (p errorProtoIO) ReadTimeout(d time.Duration) (proto.Message, error) {
	return nil, p
}
func (p errorProtoIO) Write(msg proto.Message) error {
	return p
}
func (p errorProtoIO) WriteTimeout(msg proto.Message, d time.Duration) error {
	return p
}
func (p errorProtoIO) Close() error {
	return nil
}

// Pool proto io 连接池
type Pool struct {
	ch     chan *poolProtoIO
	new    func() ProtoIO
	rw     sync.RWMutex
	closed bool
}

// NewPool 创建一个 proto io 连接池
func NewPool(newProtoIO func() (ProtoIO, error), maxCnt int) *Pool {
	p := new(Pool)
	p.ch = make(chan *poolProtoIO, maxCnt)
	p.new = func() ProtoIO {
		pio, err := newProtoIO()
		if err != nil {
			return errorProtoIO{
				error: err,
			}
		}
		return &poolProtoIO{
			pio:  pio,
			pool: p,
		}
	}
	return p
}

// Get 获取一个 proto io 对象
func (p *Pool) Get() ProtoIO {
	p.rw.RLock() // for close
	defer p.rw.RUnlock()

	if p.closed {
		return errorProtoIO{
			error: errors.New("pool closed"),
		}
	}

	select {
	case pio := <-p.ch:
		return pio
	default:
		return p.new()
	}
}

func (p *Pool) put(pio *poolProtoIO) {
	p.rw.RLock() // for close
	defer p.rw.RUnlock()

	if p.closed {
		pio.pio.Close()
		return
	}

	select {
	case p.ch <- pio:
	default:
		pio.pio.Close()
	}
}

// Close 关闭连接池, 释放所有连接
func (p *Pool) Close() {
	p.rw.Lock()
	defer p.rw.Unlock()

	p.closed = true
	close(p.ch)
	for pio := range p.ch {
		pio.pio.Close()
	}
}
