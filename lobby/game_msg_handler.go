package lobby

import (
	"sanguosha.com/sgs_herox/proto/smsg"

	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/network/msgprocessor"
	"sanguosha.com/sgs_herox"
)

func initGameMsgHandler(app *appframe.Application) {
	app.RegisterResponse((*smsg.RespResponseTest)(nil))
	app.RegisterResponse((*smsg.LoGaRespLookGame)(nil))

	appframe.ListenMsgSugar(app, ServerNotice)
	appframe.ListenMsgSugar(app, EntityClientRawMessage)
	appframe.ListenMsgSugar(app, onGaLoNtfGameStateChange)
}

func onGaLoNtfGameStateChange(sender appframe.Server, msg *smsg.NtfGameStateChange) {
	gameId := msg.GameUUID

	g, ok := gameMgr.findGame(gameId)
	if !ok || g == nil {
		return
	}
	switch msg.State {
	case smsg.NtfGameStateChange_Start:
		g.state = GameStateRunning
	case smsg.NtfGameStateChange_End:
		g.state = GameStateEnd
		//logrus.Debug("GameManager OnGameOver ", gameId, " ", msg)
		clearWhenGameOver(g)
	default:

	}
}

func clearWhenGameOver(g *Game) {
	gameId := g.gameId

	if room, ok := roomMgr.findRoomById(g.roomId); ok {
		room.clearGame()
	}

	//for _, userID := range g.userIds {
	//	u, ok := userMgr.findUser(userID)
	//	if !ok {
	//		continue
	//	}
	//
	//	if !u.isInGame() {
	//		continue
	//	}
	//
	//	if u.game.gameId != gameId {
	//		continue
	//	}
	//
	//	u.clearGame()
	//	u.quitRoomIfOffline()
	//}

	gameMgr.deleteGame(gameId)
}

func ServerNotice(sender appframe.Server, cmd *smsg.ServerNotice) {
	//logrus.Debug("lobby ServerNotice:", *cmd)

	serverList := AppInstance.GetAvailableServerIDs(sgs_herox.SvrTypeGate)
	for _, id := range serverList {
		AppInstance.GetServer(id).SendMsg(cmd)
	}
}

func EntityClientRawMessage(sender appframe.Server, msg *smsg.EnLoRawMessage) {
	u, ok := userMgr.findUser(msg.UserID)
	if !ok {
		return
	}
	if msg, e := msgprocessor.OnUnmarshal(msg.MsgId, msg.Data); e == nil {
		u.SendMsg(msg.(proto.Message))
	}
}

//func onGaLoNtfPlayerLeaveGame(sender appframe.Server, msg *smsg.GaLoNtfPlayerLeaveGame) {
//	u, ok := userMgr.findUser(msg.Userid)
//	if !ok {
//		return
//	}
//
//	if !u.isInGame() {
//		return
//	}
//
//	if u.game.gameId == msg.GameUUID {
//		u.clearGame()
//	}
//}
