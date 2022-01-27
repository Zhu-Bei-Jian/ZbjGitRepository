package client

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"time"

	"sanguosha.com/sgs_herox/proto/cmsg"
)

func listenMsg(client *Client) {
	//client.ListenMsg((*cmsg.SRespLogin)(nil), client.onRespLogin)
	//func (c *Client) onRespLogin(msg proto.Message) {}
	listenGameMsg(client)
}

func (c *Client) onConnect() {
	c.reqLogin()
}

func (c *Client) reqLogin() {
	c.SetWaitMsg((*cmsg.SRespLogin)(nil))
	c.Request(&cmsg.CReqLogin{
		Ticket:  c.ticket,
		Version: "",
		//Extra: &gamedef.ExtraAccountInfo{
		//	LoginType:  gameconf.AccountLoginTyp_ALTSimulator,
		//	DeviceType: gameconf.ClientDeviceTyp_CDTPC,
		//},
	}, func(resp *cmsg.SRespLogin, err error) {
		if err != nil {
			c.log.WithError(err).Errorf("Login failed")
			return
		}
		c.CheckWaitMsg((*cmsg.SRespLogin)(nil))
		//resp := msg.(*cmsg.SRespLogin)
		// 创角
		//if resp.ErrCode == cmsg.SRespLogin_ErrUninit {
		//	c.reqInitMyData()
		//	return
		//}
		//c.userID = resp.Userid
		//
		//// 完成登录
		if resp.ErrCode == 0 {
			onLoginFinish(resp.UserId, nil)
			//c.startTestAfterLoginSuccess()
			//c.reqMyData()
			c.SendMsg(&cmsg.CReqRoomQuickJoin{
				GameMode: gameconf.GameModeTyp(generateRandomMode()),
			})
			c.SendMsg(&cmsg.CReqRoomReady{
				Ready: true,
			})
		} else {
			c.log.WithFields(logrus.Fields{
				"errCode": resp.ErrCode,
			}).Errorf("Login failed")
			onLoginFinish(0, fmt.Errorf("Login failed"))
		}
	}, 5*time.Second)
}

func onLoginFinish(userid uint64, err error) {
	if err != nil {
		pushCol("login_err", 1)
		return
	}
	pushCol("login_success", 1)
}

func (c *Client) reqInitMyData() {
	//c.SetWaitMsg((*cmsg.SRespInitMyData)(nil))
	//c.Request(&cmsg.CReqInitMyData{
	//	Nickname: c.ticket,
	//}, func(resp *cmsg.SRespInitMyData, err error) {
	//	c.CheckWaitMsg((*cmsg.SRespInitMyData)(nil))
	//	// 完成创角
	//	if resp.ErrCode != 0 {
	//		c.log.WithFields(logrus.Fields{
	//			"errCode": resp.ErrCode,
	//			"errMsg":  resp.ErrMsg,
	//		}).Error("Init data failed")
	//		onLoginFinish(0, errors.New("Init data failed"))
	//	} else {
	//		c.log.Info("Init data succ")
	//		c.reqLogin()
	//	}
	//}, 5*time.Second)
}

func (c *Client) reqMyData() {
	c.SetWaitMsg((*cmsg.SRespMyData)(nil))
	c.Request(&cmsg.CReqMyData{}, func(resp *cmsg.SRespMyData, err error) {
		c.CheckWaitMsg((*cmsg.SRespMyData)(nil))
		//c.myData = resp

		time.AfterFunc(5*time.Second, func() {
			if TestConf.TestEntity {
				c.GoTestEntity()
			} else if TestConf.TestGame {
				c.GoTestGame()
			}
		})
	}, 5*time.Second)

}

func (c *Client) testReqGMCommand(cmd string) {
	cm := (*cmsg.SReqGMCommand)(nil)
	if !c.IsWaitMsg(cm) {
		c.SetWaitMsg(cm)
		c.Request(&cmsg.CReqGMCommand{
			Cmd: cmd,
		}, func(msg *cmsg.SReqGMCommand, err error) {
			c.CheckWaitMsg(msg)
		}, 5*time.Second)
	}
}

func (c *Client) reqPing(svr int32) {
	c.Request(&cmsg.CReqPing{
		TimeTag: time.Now().UnixNano(),
		SvrType: svr,
	}, func(resp *cmsg.SRespPing, err error) {
		if err != nil {
			return
		}
		diff := (time.Now().UnixNano() - resp.TimeTag)
		if diff >= 0 {
			pushCol("ping", resp.SvrType, diff)
		}
	}, 5*time.Second)
}
