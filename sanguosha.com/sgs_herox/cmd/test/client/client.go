package client

import (
	"errors"
	"os"
	"reflect"
	"sanguosha.com/baselib/log"
	"sanguosha.com/sgs_herox/gameutil"
	"sync"
	"sync/atomic"
	"time"

	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"net"
	"sanguosha.com/baselib/appframe/pio"
	"sanguosha.com/baselib/framework/netframe"
	"sanguosha.com/baselib/framework/worker"
	"sanguosha.com/sgs_herox/proto/cmsg"
)

var (
	// ErrTimeout ...
	ErrTimeout = errors.New("ErrTimeout")
)

var loggerCloser func()

func Init(isStressTesting bool) {
	level := logrus.DebugLevel
	if isStressTesting {
		level = logrus.InfoLevel
	}

	loggerCloser, _ = log.InitLogrus(&log.Config{
		Name:  "client",
		Level: int(level),
		//UseJSON: true,
		Outputs: map[string]map[string]interface{}{
			"file": {
				"path":   "./logs",
				"rotate": true,
				"json":   false,
			},
		},
	})
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		//DisableTimestamp: true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableSorting:  true,
	})
	logrus.SetOutput(os.Stdout)
}

// MsgHandler ...
type MsgHandler func(msg proto.Message)

type respHandler func(msg proto.Message, err error)

// Client ...
type Client struct {
	testPrefix string
	testMax    int

	svrAddr string
	ticket  string
	log     logrus.FieldLogger

	UserData

	pio    pio.ProtoIO
	ch     chan func()
	closed int32

	msgHandlers map[reflect.Type]MsgHandler
	waitResp    map[reflect.Type][]respHandler
	mtx         sync.RWMutex
}

// New ...
func New(svrAddr string, ticket string, prefix string, max int) *Client {
	return &Client{
		testPrefix: prefix,
		testMax:    max,

		svrAddr:     svrAddr,
		ticket:      ticket,
		log:         logrus.WithField("ticket", ticket),
		ch:          make(chan func(), 256),
		msgHandlers: map[reflect.Type]MsgHandler{},
		waitResp:    map[reflect.Type][]respHandler{},
	}
}

// Logger ...
func (c *Client) Logger() logrus.FieldLogger {
	return c.log
}

// ListenMsg ...
func (c *Client) ListenMsg(msg proto.Message, handler MsgHandler) {
	regMsg(msg)
	typ := reflect.TypeOf(msg)

	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.msgHandlers[typ] = handler
}

// SendMsg ...
func (c *Client) SendMsg(msg proto.Message) error {
	regMsg(msg)
	return c.pio.Write(msg)
}

// Request 由于没有 seqid, 因此测试客户端的 Request 会有些使用限制, 使用时需要注意.
func (c *Client) Request(msg proto.Message, respCbk interface{}, timeout time.Duration) {
	v := reflect.ValueOf(respCbk)

	// type check.
	if v.Type().NumIn() != 2 {
		logrus.Panic("Client Request respCbk params num wrong")
	}
	var tempError error
	if v.Type().In(1) != reflect.TypeOf(&tempError).Elem() {
		logrus.Panic("Client Request respCbk params num in 1 is not a error")
	}

	respType := v.Type().In(0)
	regMsg(reflect.New(respType).Elem().Interface().(proto.Message))

	popCbk := func() (respHandler, bool) {
		var cbk respHandler

		c.mtx.Lock()
		cbks := c.waitResp[respType]
		if len(cbks) == 0 {
			c.mtx.Unlock()
			return nil, false
		}
		cbk, c.waitResp[respType] = cbks[0], cbks[1:]
		c.mtx.Unlock()

		return cbk, true
	}

	c.mtx.Lock()
	c.msgHandlers[respType] = func(msg proto.Message) {
		if cbk, ok := popCbk(); ok {
			cbk(msg, nil)
		}
	}
	cancelTimer := worker.AfterPost(c, timeout, func() {
		if cbk, ok := popCbk(); ok {
			cbk(nil, ErrTimeout)
		}
	})
	c.waitResp[respType] = append(c.waitResp[respType], func(msg proto.Message, err error) {
		cancelTimer()
		var in [2]reflect.Value
		if msg == nil {
			in[0] = reflect.Zero(respType)
		} else {
			in[0] = reflect.ValueOf(msg)
		}
		if err == nil {
			in[1] = reflect.Zero(reflect.TypeOf(&err).Elem())
		} else {
			in[1] = reflect.ValueOf(err)
		}
		v.Call(in[:])
	})
	c.mtx.Unlock()

	err := c.SendMsg(msg)
	if err != nil {
		worker.AfterPost(c, 0, func() {
			if cbk, ok := popCbk(); ok {
				cbk(nil, err)
			}
		})
	}
}

// Post ...
func (c *Client) Post(f func()) {
	if atomic.LoadInt32(&c.closed) == 0 {
		c.ch <- f
	}
}

// Run ...
func (c *Client) Run() (err error) {

	listenMsg(c)

	conn, err := net.Dial("tcp", c.svrAddr)
	if err != nil {
		c.log.WithError(err).WithField("addr", c.svrAddr).Error("New pio error")
		onLoginFinish(0, errors.New("New pio failed"))
		return err
	}
	c.pio = pio.New(conn, msg2id, id2msg, binary.LittleEndian, nil)
	defer c.pio.Close()

	go func() {
		for {
			msg, err := c.pio.Read()
			if atomic.LoadInt32(&c.closed) != 0 {
				return
			}
			if err != nil {
				// skip unregister msg
				if err == pio.ErrID2Msg {
					continue
				}
				c.log.WithError(err).Error("pio read error")
				c.Exit()
				return
			}

			typ := reflect.TypeOf(msg)
			c.mtx.RLock()
			f, ok := c.msgHandlers[typ]
			c.mtx.RUnlock()
			if ok {
				c.Post(func() {
					f(msg)
				})
			}
		}
	}()

	// 心跳维护.
	heartBeatTicker := time.NewTicker(30 * time.Second)
	defer heartBeatTicker.Stop()
	go func() {
		for range heartBeatTicker.C {
			c.SendMsg(&netframe.Client_ReqHeartBeat{TimeStamp: time.Now().Unix()})
		}
	}()

	c.ListenMsg((*cmsg.SNoticeLogout)(nil), func(msg proto.Message) {
		c.log.Warn("Server notice logout")
		c.ClearGame()
	})

	c.onConnect()

	for f := range c.ch {
		gameutil.SafeCall(f)
	}

	return nil
}

// Exit ...
func (c *Client) Exit() {
	if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		close(c.ch)
	}
}
