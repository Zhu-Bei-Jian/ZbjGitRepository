package lobby

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
)

func ListenServerEvent(app *appframe.Application) {
	app.ListenServerEvent(sgs_herox.SvrTypeGame, onGameServerEvent)
	app.ListenServerEvent(sgs_herox.SvrTypeAI, onAIServerEvent)
	app.ListenServerEvent(sgs_herox.SvrTypeEntity, onEntityServerEvent)
	app.ListenServerEvent(sgs_herox.SvrTypeGate, onGateServerEvent)
}

func onGameServerEvent(svrid uint32, e netcluster.SvrEvent) {
	switch e {
	case netcluster.SvrEventStart, netcluster.SvrEventReconnect:
	case netcluster.SvrEventQuit, netcluster.SvrEventDisconnect:
		logrus.Error(fmt.Sprintf("SvrTypeGame %d Quit or Disconnect", svrid))
		for _, g := range gameMgr.games {
			if svrid != g.node.ID() {
				continue
			}
			clearWhenGameOver(g)
		}
	}
}

func onAIServerEvent(svrid uint32, e netcluster.SvrEvent) {
	switch e {
	case netcluster.SvrEventStart, netcluster.SvrEventReconnect:
	case netcluster.SvrEventQuit, netcluster.SvrEventDisconnect:
		logrus.Error(fmt.Sprintf("SvrTypeAI %d Quit or Disconnect", svrid))
		for _, g := range gameMgr.games {
			if g.aiSvrID == svrid {
				g.aiSvrID = 0
			}
		}
	}
}

func onEntityServerEvent(svrID uint32, event netcluster.SvrEvent) {
	switch event {
	case netcluster.SvrEventQuit, netcluster.SvrEventDisconnect:
		gate2userIds := make(map[uint32][]uint64)
		userMgr.execByEverySession(func(u *user, isLast bool) {
			entityID := AppCfg.GetUserEntityID(u.userid)
			if entityID != svrID {
				return
			}

			gateId := u.session.SvrID
			if gate2userIds[gateId] == nil {
				gate2userIds[gateId] = make([]uint64, 0)
			}
			gate2userIds[gateId] = append(gate2userIds[gateId], u.userid)

			if !isLast {
				return
			}

			workerMgr.Post(func() {
				for gateId, userIds := range gate2userIds {
					AppInstance.GetServer(gateId).SendMsg(&smsg.AllGaNtfKickUserOut{
						KickAll: false,
						UserIds: userIds,
						Reason:  gameconf.KickUserOutReason_KUORelogin,
					})
					for _, userId := range userIds {
						u, exist := userMgr.findUser(userId)
						if !exist {
							continue
						}
						u.onDisconnect()
					}
				}
			})

		})
	}
}

func onGateServerEvent(svrID uint32, event netcluster.SvrEvent) {
	switch event {
	case netcluster.SvrEventQuit, netcluster.SvrEventDisconnect:
		userIds := make([]uint64, 0)
		userMgr.execByEverySession(func(u *user, isLast bool) {
			gateId := u.session.SvrID
			if gateId != svrID {
				return
			}

			userIds = append(userIds, u.userid)

			if !isLast {
				return
			}

			workerMgr.Post(func() {
				for _, userId := range userIds {
					u, exist := userMgr.findUser(userId)
					if !exist {
						continue
					}
					u.onDisconnect()
				}
			})
		})
	}
}
