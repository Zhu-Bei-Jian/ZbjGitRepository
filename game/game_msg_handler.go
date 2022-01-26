package game

import (
	"github.com/golang/protobuf/proto"
	"reflect"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared/notifier"

	"sanguosha.com/sgs_herox/proto/cmsg"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"

	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
)

func InitGameMsgHandler(app *appframe.Application) {
	logrus.Info("InitGameMsgHandler !!!!!!!!!!!")

	app.RegisterResponse((*smsg.RespUserData)(nil))
	app.RegisterResponse((*smsg.PuEnRespAddProp)(nil))

	appframe.ListenRequestSugar(app, LoGaReqNewGame)
	appframe.ListenRequestSugar(app, LoGaReqUserQuit)

	appframe.ListenRequestSugar(app, LoGaReqLookGame)
	appframe.ListenRequestSugar(app, onReqGameRemove)

	appframe.ListenMsgSugar(app, UserDisconnect)
	appframe.ListenMsgSugar(app, UserConnect)

	appframe.ListenSessionMsgSugar(app, onReqGameScene)
	appframe.ListenSessionMsgSugar(app, onReqChat)

	ListenGameMsgSugar(app, onReqAct)
	ListenGameMsgSugar(app, onReqOpt)
	ListenGameMsgSugar(app, onReqCancel)
}

func LoGaReqLookGame(sender appframe.Requester, msg *smsg.LoGaReqLookGame) {
	gameId := msg.GameUUID

	resp := &smsg.LoGaRespLookGame{
		GameUUID: gameId,
	}
	defer sender.Resp(resp)

	g, exist := gameMgr.findGame(gameId)
	if !exist {
		resp.ErrCode = smsg.LoGaRespLookGame_ErrLookGameFailed
		return
	}

	sessionId := appframe.SessionID{
		SvrID: msg.Looker.GateId,
		ID:    msg.Looker.Session,
	}
	session := App.GetSession(sessionId)
	u := newUser(msg.Looker.UserId, msg.Looker.SeatId, msg.Looker.UserBrief, session)
	player := g.onLookerJoin(u)
	playerMgr.add(player)
}

func LoGaReqNewGame(sender appframe.Requester, msg *smsg.LoGaReqNewGame) {
	gameId := msg.GameUUID

	resp := &smsg.LoGaRespNewPVPGame{
		GameUUID: gameId,
	}
	defer sender.Resp(resp)

	game := newGame(App, msg.RoomSetting, gameId, msg.RoomId, msg.RoomNO, gCfg)
	g := game.base
	err := gameMgr.add(gameId, g)
	if err != nil {
		resp.ErrCode = smsg.LoGaRespNewPVPGame_ErrCreateGameFailed
		return
	}

	registerStartEndCallback(g)
	for _, v := range msg.Users {
		sessionId := appframe.SessionID{
			SvrID: v.GateId,
			ID:    v.Session,
		}
		session := App.GetSession(sessionId)
		u := newUser(v.UserId, v.SeatId, v.UserBrief, session)
		player := g.OnUserJoin(u, v.SeatId)

		playerMgr.add(player)
	}

	for _, v := range msg.Lookers {
		sessionId := appframe.SessionID{
			SvrID: v.GateId,
			ID:    v.Session,
		}
		session := App.GetSession(sessionId)
		u := newUser(v.UserId, v.SeatId, v.UserBrief, session)
		player := g.onLookerJoin(u)

		playerMgr.add(player)
	}

	g.prepareStart()
}

func registerStartEndCallback(g *GameBase) {
	gameId := g.gameUUID
	g.startCallback = func() {
		msg := &smsg.NtfGameStateChange{
			GameUUID: gameId,
			State:    smsg.NtfGameStateChange_Start,
		}
		App.GetService(sgs_herox.SvrTypeLobby).SendMsg(msg)
	}

	g.clearCallback = func() {
		game, exist := gameMgr.findGame(gameId)
		if !exist {
			return
		}
		if game != g {
			return
		}
		for _, v := range game.players {
			if v == nil {
				continue
			}
			playerMgr.del(v)
		}

		for _, v := range game.lookers {
			if v == nil {
				continue
			}
			playerMgr.del(v)
		}

		gameMgr.deleteGame(gameId)
		msg := &smsg.NtfGameStateChange{
			GameUUID: gameId,
			State:    smsg.NtfGameStateChange_End,
		}
		App.GetService(sgs_herox.SvrTypeLobby).SendMsg(msg)
	}
}

func LoGaReqUserQuit(sender appframe.Requester, msg *smsg.LoGaReqUserQuit) {
	userId := msg.Userid
	resp := &smsg.LoGaRespUserQuit{
		Userid: userId,
	}
	defer sender.Resp(resp)

	player, exist := playerMgr.findPlayerByUserId(userId)
	if !exist {
		return
	}

	oldSessionId := player.GetUser().SessionID()
	playerMgr.delSession(oldSessionId)

	g := player.game

	if g.isLooker(userId) {
		playerMgr.del(player)
		g.onLookerQuit(player)
	} else {
		g.onPlayerQuit(player)
		if player.user.userBrief.AccountType == gameconf.AccountLoginTyp_ALTTablePark {
			tpNotifier.NotifyGameStatusAsync(int32(g.roomNO), player.user.userBrief.ThirdAccountId, notifier.TableParkGameStatus_Over)
		}
	}
}

func UserDisconnect(sender appframe.Server, msg *smsg.UserDisconnect) {
	userId := msg.Userid

	player, exist := playerMgr.findPlayerByUserId(userId)
	if !exist {
		return
	}

	oldSessionId := player.GetUser().SessionID()
	playerMgr.delSession(oldSessionId)

	player.game.onPlayerDisconnect(player)
}

func UserConnect(sender appframe.Server, msg *smsg.UserConnect) {
	userId := msg.Userid

	player, exist := playerMgr.findPlayerByUserId(userId)
	if !exist {
		return
	}

	oldSessionId := player.GetUser().SessionID()

	sessionId := appframe.SessionID{
		SvrID: msg.Gateid,
		ID:    msg.Session,
	}
	session := App.GetSession(sessionId)
	player.game.onPlayerReconnect(player, session)

	playerMgr.delSession(oldSessionId)
	playerMgr.addSession(sessionId, player)
}

func ListenGameMsgSugar(app *appframe.Application, msgHandler interface{}) {
	v := reflect.ValueOf(msgHandler)

	// type check.
	if v.Type().NumIn() != 2 {
		logrus.Panic("ListenSessionMsgSugar handler params num wrong")
	}

	var tempPlayer *Player
	if v.Type().In(0) != reflect.TypeOf(&tempPlayer).Elem() {
		logrus.Panic("ListenSessionMsgSugar handler num in 0 is not Session")
	}

	iMsg := reflect.New(v.Type().In(1)).Elem().Interface()
	msg := iMsg.(proto.Message)
	app.ListenSessionMsg(msg, func(sender appframe.Session, msg proto.Message) {
		senderID := sender.ID()
		p, exist := playerMgr.findPlayerBySessionID(senderID)
		if !exist {
			resp := &cmsg.SRespCommonError{ErrCode: cmsg.SRespCommonError_SessionNotInGame}
			sender.SendMsg(resp)
			return
		}

		//旁观玩家不可以发送游戏内消息
		if p.seatId < 0 {
			resp := &cmsg.SRespCommonError{ErrCode: cmsg.SRespCommonError_LookerCannotOp}
			sender.SendMsg(resp)
			return
		}
		g := p.GetGame()
		g.DoNow(func() {
			g.CheckAllow(p, msg, func() {
				v.Call([]reflect.Value{reflect.ValueOf(p), reflect.ValueOf(msg)})
				//g.Print()
			})
		})
	})
}

func onReqAct(p *Player, msg *cmsg.CReqAct) {
	g := p.GetGame()
	g.onReqAct(p, msg)
}

func onReqCancel(p *Player, msg *cmsg.CReqCancelCurOpt) {
	g := p.GetGame()
	g.onReqCancelCurOpt(p, msg)
}

func onReqOpt(p *Player, msg *cmsg.CReqOpt) {
	a := p.GetGame().ActDataPeek()
	if a == nil {
		return
	}
	switch t := a.(type) {
	case *ActionSelectCard:
		t.OnMessage(p, msg)
	case *ActionSelectCamp:
		t.OnMessage(p, msg)
	default:

	}
}

func onReqChat(sender appframe.Session, msg *cmsg.CReqChat) {
	resp := &cmsg.SRespChat{}
	p, exist := playerMgr.findPlayerBySessionID(sender.ID())
	if !exist {
		resp.ErrCode = cmsg.SRespChat_ErrNotInGame
		sender.SendMsg(resp)
		return
	}

	g := p.GetGame()
	g.onReqChat(p, msg)
}

func onReqGameScene(sender appframe.Session, req *cmsg.CReqGameScene) {
	resp := &cmsg.SRespGameScene{}
	p, exist := playerMgr.findPlayerBySessionID(sender.ID())
	if !exist {
		resp.ErrCode = cmsg.SRespGameScene_ErrGameNotFound
		sender.SendMsg(resp)
		return
	}

	if p.game.IsOver() {
		resp.ErrCode = cmsg.SRespGameScene_ErrGameOver
		sender.SendMsg(resp)
		return
	}
	p.game.onReqGameScene(p, req)
}

func onReqGameRemove(sender appframe.Requester, req *smsg.ReqGameRemove) {
	resp := &smsg.RespGameRemove{}
	defer sender.Resp(resp)
	if len(req.GameIDs) != 0 {
		for _, gameID := range req.GameIDs {
			g, ok := gameMgr.findGame(gameID)
			if !ok {
				continue
			}
			if g.IsOver() {
				continue
			}
			resp.GameIDs = append(resp.GameIDs, gameID)
			g.DoNow(func() {
				if g.IsOver() {
					return
				}
				g.ForceOver()
			})
		}
		return
	}
	allGames := gameMgr.allGames()
	for _, game := range allGames {
		g := game
		if g.IsOver() {
			continue
		}
		resp.GameIDs = append(resp.GameIDs, g.gameUUID)
		g.ForceOver()
	}
}
