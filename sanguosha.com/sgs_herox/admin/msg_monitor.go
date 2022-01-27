package admin

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/mq"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/smsg"

	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type MsgMonitor struct {
	connId int32

	consumer mq.Consumer
	clients  sync.Map

	workingCMD atomic.Value
}

func NewMsgMonitor(mqNodes map[string]*config.MQNode) (*MsgMonitor, error) {
	mqCfg, exist := mqNodes[config.MQNode_Monitor]
	if !exist {
		return nil, errors.New("no MQNode_Monitor")
	}

	consumer, err := mq.NewConsumer(mq.Config{
		Open:    mqCfg.Open,
		Type:    mqCfg.Type,
		Address: mqCfg.Address,
	})

	if err != nil {
		return nil, err
	}

	p := &MsgMonitor{
		consumer: consumer,
	}
	p.workingCMD.Store("")

	err = p.consumer.Sub(config.TopicMonitor, p.subCallback)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *MsgMonitor) subCallback(data []byte) {
	p.clients.Range(func(key, value interface{}) bool {
		c := key.(*MClient)
		c.sendByte(data)
		c.sendStr("---------------------------------------------------------------")
		return true
	})

}

func (p *MsgMonitor) AddConn(conn *websocket.Conn) {
	client := &MClient{id: "tempid", socket: conn, send: make(chan []byte, 1024), m: p}

	p.clients.Store(client, true)

	util.SafeGo(func() { client.read() })
	util.SafeGo(func() { client.write() })
}

func (p *MsgMonitor) gmAll(typ smsg.AsAllReqMSGMonitor_Type, userId uint64, c *MClient) {
	var wg sync.WaitGroup
	allServerIds := make([]uint32, 0)
	successServerIds := make([]uint32, 0)
	failServerIds := make([]uint32, 0)

	reqMsg := smsg.AsAllReqMSGMonitor{
		Type:   typ,
		UserId: userId,
	}

	for svrType := sgs_herox.SvrTypeGate; svrType < sgs_herox.SvrTypeEnd; svrType++ {
		if svrType == sgs_herox.SvrTypeAdmin || svrType == sgs_herox.SvrTypeAI {
			continue
		}
		serverIds := app.GetAvailableServerIDs(svrType)
		for _, v := range serverIds {
			wg.Add(1)
			serverID := v
			allServerIds = append(allServerIds, serverID)
			app.GetServer(serverID).ReqSugar(&reqMsg, func(resp *smsg.AsAllRespMSGMonitor, err error) {
				defer wg.Done()
				if err != nil {
					failServerIds = append(failServerIds, serverID)
					return
				}

				if resp.ErrCode != smsg.AsAllRespMSGMonitor_Invalid {
					failServerIds = append(failServerIds, serverID)
					return
				}
				successServerIds = append(successServerIds, serverID)
			}, time.Second*60)
		}
	}
	wg.Wait()

	if c != nil {
		c.sendStr(fmt.Sprintf("success:%v fail:%v\n", successServerIds, failServerIds))
	}
}

func (p *MsgMonitor) gmGate(typ smsg.AsAllReqMSGMonitor_Type, userId uint64, slowresponseSlowerThan int32, c *MClient) {
	var wg sync.WaitGroup
	allServerIds := make([]uint32, 0)
	successServerIds := make([]uint32, 0)
	failServerIds := make([]uint32, 0)

	reqMsg := smsg.AsAllReqMSGMonitor{
		Type:       typ,
		UserId:     userId,
		SlowerThan: slowresponseSlowerThan,
	}

	serverIds := app.GetAvailableServerIDs(sgs_herox.SvrTypeGate)
	for _, v := range serverIds {
		wg.Add(1)
		serverID := v
		allServerIds = append(allServerIds, serverID)
		app.GetServer(serverID).ReqSugar(&reqMsg, func(resp *smsg.AsAllRespMSGMonitor, err error) {
			defer wg.Done()
			if err != nil {
				failServerIds = append(failServerIds, serverID)
				return
			}

			if resp.ErrCode != smsg.AsAllRespMSGMonitor_Invalid {
				failServerIds = append(failServerIds, serverID)
				return
			}
			successServerIds = append(successServerIds, serverID)
		}, time.Second*60)
	}

	wg.Wait()

	if c != nil {
		c.sendStr(fmt.Sprintf("success:%v fail:%v\n", successServerIds, failServerIds))
	}
}

type MClient struct {
	id         string
	socket     *websocket.Conn
	send       chan []byte
	m          *MsgMonitor
	workingCMD string
}

func (c *MClient) setWorkingCMD(cmd string) {
	c.workingCMD = cmd
	c.m.workingCMD.Store(cmd)
}

func (c *MClient) stopWorkingCMD() {
	switch c.workingCMD {
	case "trace":
		c.m.gmAll(smsg.AsAllReqMSGMonitor_All_StopPrint, 0, c)
		c.m.workingCMD.Store("")
	case "slowresponse", "inout":
		c.m.gmGate(smsg.AsAllReqMSGMonitor_All_StopPrint, 0, 0, c)
		c.m.workingCMD.Store("")
	default:

	}

	c.workingCMD = ""
}

func (c *MClient) write() {
	defer c.socket.Close()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func (c *MClient) read() {
	defer func() {
		c.m.clients.Delete(c)
		c.socket.Close()
	}()

	defer c.stopWorkingCMD()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		args := strings.Split(string(msg), " ")
		if len(args) == 0 {
			c.sendStr("请输入有效命令")
			continue
		}

		cmd := args[0]

		workingCMD := c.m.workingCMD.Load().(string)
		if (cmd == "trace" || cmd == "inout" || cmd == "slowresponse") && workingCMD != "" {
			c.sendStr(fmt.Sprintf("%s 工作中,请关闭上个命令后再试", workingCMD))
			continue
		}

		switch strings.ToLower(cmd) {
		case "trace":
			if len(args) != 2 {
				c.sendStr("此命令需要参数，请重新输入")
				continue
			}

			userId, err := gameutil.ToUInt64(args[1])
			if err != nil {
				c.sendStr("参数不合法，请重新输入")
				continue
			}
			c.setWorkingCMD(cmd)
			c.m.gmAll(smsg.AsAllReqMSGMonitor_All_PrintUserMsgDetail, userId, c)
		case "slowresponse":
			if len(args) != 2 {
				c.sendStr("此命令需要参数，请重新输入")
				continue
			}

			ms, err := gameutil.ToInt32(args[1])
			if err != nil {
				c.sendStr("参数不合法，请重新输入")
				continue
			}

			if !Develop {
				if ms < 100 {
					c.sendStr("参数值不可低于100，请重新输入")
					continue
				}
			}

			c.setWorkingCMD(cmd)
			c.m.gmGate(smsg.AsAllReqMSGMonitor_Gate_PrintSlowResponse, 0, ms, c)
		case "inout":
			if len(args) != 2 {
				c.sendStr("此命令需要参数，请重新输入")
				continue
			}

			userId, err := gameutil.ToUInt64(args[1])
			if err != nil {
				c.sendStr("参数不合法，请重新输入")
				continue
			}
			c.setWorkingCMD(cmd)
			c.m.gmGate(smsg.AsAllReqMSGMonitor_Gate_PrintUserMsgInOut, userId, 0, c)
		case "stop":
			c.stopWorkingCMD()
		case "forcestop":
			c.m.gmAll(smsg.AsAllReqMSGMonitor_All_StopPrint, 0, c)
			c.setWorkingCMD("")
		default:
			c.sendStr("不支持的命令,请重新输入")
			continue
		}
	}
}

func (c *MClient) sendStr(data string) {
	c.send <- []byte(data)
}

func (c *MClient) sendByte(data []byte) {
	c.send <- data
}

func MsgMonitorPage(res http.ResponseWriter, req *http.Request, l *LoginInfo) {
	HttpWrite(res, "web/template/msg_monitor.html", nil)
}

func WsMsgMonitorConnect(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	msgMonitor.AddConn(conn)
}
