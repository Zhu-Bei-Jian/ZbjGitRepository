package client

import (
	"sanguosha.com/sgs_herox"
	"time"
)

func (c *Client) GoTestEntity() {
	tickerReq := time.NewTicker(time.Second / time.Duration(TestConf.ReqPerSec))
	tickerPing := time.NewTicker(time.Duration(TestConf.PingDelay) * time.Millisecond)
	go func() {
		for {
			select {
			case <-tickerReq.C:
				c.Post(func() {
					c.testReqGMCommand("addItem 301020 1")
					return
				})
			case <-tickerPing.C:
				c.Post(func() {
					c.reqPing(int32(sgs_herox.SvrTypeEntity))
					return
				})
			}
		}
	}()
}
