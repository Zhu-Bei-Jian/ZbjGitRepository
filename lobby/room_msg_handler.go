package lobby

import (
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"reflect"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/proto/cmsg"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
)

func initRoomMsgHandler(app *appframe.Application) {
	app.RegisterResponse((*smsg.LoGaRespUserQuit)(nil))

	appframe.ListenSessionMsgSugar(app, onReqMyRoom)
	appframe.ListenSessionMsgSugar(app, onReqRoomCreate)
	appframe.ListenSessionMsgSugar(app, onReqRoomStartGame)

	ListenSessionMsgSugarUser(app, onReqRoomSettingChange)
	ListenSessionMsgSugarUser(app, onReqRoomJoin)
	ListenSessionMsgSugarUser(app, onReqRoomSeatChange)
	ListenSessionMsgSugarUser(app, onReqRoomLookerInfo)
	ListenSessionMsgSugarUser(app, onReqRoomLeave)
	ListenSessionMsgSugarUser(app, onReqRoomKick)
	ListenSessionMsgSugarUser(app, onReqRoomReady)
	ListenSessionMsgSugarUser(app, onReqRoomQuickJoin)
}

func ListenSessionMsgSugarUser(app *appframe.Application, msgHandler interface{}) {
	v := reflect.ValueOf(msgHandler)

	// type check.
	if v.Type().NumIn() != 2 {
		logrus.Panic("ListenSessionMsgSugar handler params num wrong")
	}

	var tempuser *user
	if v.Type().In(0) != reflect.TypeOf(&tempuser).Elem() {
		logrus.Panic("ListenSessionMsgSugar handler num in 0 is not user")
	}

	iMsg := reflect.New(v.Type().In(1)).Elem().Interface()
	msg := iMsg.(proto.Message)
	app.ListenSessionMsg(msg, func(sender appframe.Session, msg proto.Message) {
		senderID := sender.ID()
		u, exist := userMgr.findUserBySessionID(senderID)
		if !exist {
			resp := &cmsg.SRespCommonError{ErrCode: cmsg.SRespCommonError_SessionNotInLobby}
			sender.SendMsg(resp)
			return
		}

		v.Call([]reflect.Value{reflect.ValueOf(u), reflect.ValueOf(msg)})
	})
}

func onReqMyRoom(sender appframe.Session, req *cmsg.CReqMyRoom) {
	u, ok := userMgr.findUserBySessionID(sender.ID())
	if !ok {
		return
	}

	resp := &cmsg.SRespMyRoom{}
	defer sender.SendMsg(resp)

	room := u.room
	if room != nil {
		resp.Room = room.ToCMSG()
	}
}

func onReqRoomCreate(sender appframe.Session, req *cmsg.CReqRoomCreate) {
	u, ok := userMgr.findUserBySessionID(sender.ID())
	if !ok {
		return
	}

	resp := &cmsg.SRespRoomCreate{}
	defer sender.SendMsg(resp)

	if u.isInRoom() {
		resp.ErrCode = cmsg.SRespRoomCreate_ErrAlreadyInRoom
		return
	}

	if u.isInGame() {
		resp.ErrCode = cmsg.SRespRoomCreate_ErrAlreadyInGame
		return
	}

	setting := req.Settting

	err := roomMgr.checkRoomCfg(setting)
	if err != nil {
		resp.ErrCode = cmsg.SRespRoomCreate_ErrRoomSetting
		return
	}

	room, err := roomMgr.newRoom(setting, u.userid)
	if err != nil {
		resp.ErrCode = cmsg.SRespRoomCreate_ErrCreateRoom
		return
	}

	_, ok = room.Enter(u, true)
	if !ok {
		resp.ErrCode = cmsg.SRespRoomCreate_ErrSitDownFail
		return
	}

	resp.Room = room.ToCMSG()
}

//请求快速进房间，如无可用房间，将创建新房间
func onReqRoomQuickJoin(u *user, req *cmsg.CReqRoomQuickJoin) {
	resp := &cmsg.SRespRoomQuickJoin{}
	defer u.SendMsg(resp)

	gameMode := req.GameMode

	//默认文字场
	if gameMode == gameconf.GameModeTyp_MGTInvalid {
		gameMode = gameconf.GameModeTyp_MGTSpyText
	}

	if u.isInGame() {
		resp.ErrCode = cmsg.SRespRoomQuickJoin_ErrAlreadyInGame
		return
	}

	if u.isInRoom() {
		resp.Room = u.room.ToCMSG()
		return
	}

	var err error

	room, ok := roomMgr.findRoomCanEnter(gameMode)
	if !ok {
		room, err = roomMgr.newRoom(roomMgr.defaultRoomSetting(gameMode), u.userid)
		if err != nil {
			resp.ErrCode = cmsg.SRespRoomQuickJoin_ErrCreateRoom
			return
		}
	}
	_, ok = room.Enter(u, true)
	if !ok {
		resp.ErrCode = cmsg.SRespRoomQuickJoin_ErrSitDownFail
		return
	}

	resp.Room = room.ToCMSG()
}

func onReqRoomSettingChange(u *user, req *cmsg.CReqRoomSettingChange) {
	resp := &cmsg.SRespRoomSettingChange{}
	defer u.SendMsg(resp)

	if u.isInGame() {
		resp.ErrCode = cmsg.SRespRoomSettingChange_ErrAlreadyInGame
		return
	}

	if !u.isInRoom() {
		resp.ErrCode = cmsg.SRespRoomSettingChange_ErrNotInRoom
		return
	}

	room := u.room

	if !room.isOwner(u.userid) {
		resp.ErrCode = cmsg.SRespRoomSettingChange_ErrNoAuthority
		return
	}

	setting := req.Setting

	err := roomMgr.checkRoomCfg(setting)
	if err != nil {
		resp.ErrCode = cmsg.SRespRoomSettingChange_ErrRoomSetting
		return
	}

	seatCount := room.SeatCount()

	//房间人数有改变
	if uint32(seatCount) != setting.MaxPlayer {
		ok := room.changeSeatCount(setting.MaxPlayer)
		if !ok {
			resp.ErrCode = cmsg.SRespRoomSettingChange_ErrChangeSeatCount
			return
		}
	}

	room.updateSetting(setting, u.userid)

	resp.Setting = setting
}

func onReqRoomJoin(u *user, req *cmsg.CReqRoomJoin) {
	roomNO := req.RoomNO

	resp := &cmsg.SRespRoomJoin{}
	defer u.SendMsg(resp)

	room, exist := roomMgr.findRoomByRoomNO(roomNO)
	if !exist {
		resp.ErrCode = cmsg.SRespRoomJoin_ErrRoomNotExist
		return
	}

	if !room.isAllowEnter() {
		resp.ErrCode = cmsg.SRespRoomJoin_ErrRoomNotAllowEnter
		return
	}

	if room.isInGame() {
		if room.isLooker(u.userid) {
			resp.ErrCode = cmsg.SRespRoomJoin_ErrAlreadyLooker
			return
		}

		if u.isInGame() {
			resp.ErrCode = cmsg.SRespRoomJoin_ErrPlayerInGame
			return
		}

		_, ok := room.Enter(u, false)
		if !ok {
			resp.ErrCode = cmsg.SRespRoomJoin_ErrSitDownFail
			return
		}

		room.lookGame(u)
		resp.RoomNO = roomNO
		resp.RoomIsInGame = true
		resp.IsLooker = true
		resp.Room = room.ToCMSG()
		return
	}

	if room.isSeatFull() {
		_, ok := room.Enter(u, false)
		if !ok {
			resp.ErrCode = cmsg.SRespRoomJoin_ErrSitLookerFailed
			return
		}
		resp.RoomNO = roomNO
		resp.IsLooker = true
		resp.Room = room.ToCMSG()
		return
	}

	_, ok := room.Enter(u, true)
	if !ok {
		resp.ErrCode = cmsg.SRespRoomJoin_ErrSitDownFail
		return
	}
	resp.RoomNO = roomNO
	resp.Room = room.ToCMSG()
}

func onReqRoomSeatChange(u *user, req *cmsg.CReqRoomSeatChange) {
	targetSeatId := req.TargetSeatId

	resp := &cmsg.SRespRoomSeatChange{}
	defer u.SendMsg(resp)

	if !u.isInRoom() {
		resp.ErrCode = cmsg.SRespRoomSeatChange_ErrNotInRoom
		return
	}
	room := u.room

	if targetSeatId == -1 {
		//请求成为旁观者
		if room.isLooker(u.userid) {
			resp.ErrCode = cmsg.SRespRoomSeatChange_ErrSitUpFail_AlreadySitUp
			return
		}

		if room.getUserCount() == 1 {
			resp.ErrCode = cmsg.SRespRoomSeatChange_ErrSitUpFail_IsLastUser
			return
		}

		seatId, seat, exist := room.findSeatByUserId(u.userid)
		if !exist {
			resp.ErrCode = cmsg.SRespRoomSeatChange_ErrUnknown
			return
		}

		room.SitUpToLooker(seatId, seat)
		resp.LookerType = gameconf.LookerTyp_LTBlind
	} else {
		//请求坐下
		if !room.isLooker(u.userid) {
			resp.ErrCode = cmsg.SRespRoomSeatChange_ErrSitDownFail_AlreadySitDown
			return
		}

		seat, ok := room.findSeatBySeatId(targetSeatId)
		if !ok {
			resp.ErrCode = cmsg.SRespRoomSeatChange_ErrTargetSeatNotExist
			return
		}

		if !seat.isEmpty() {
			resp.ErrCode = cmsg.SRespRoomSeatChange_ErrTargetSeatAlreadySeated
			return
		}

		room.LookerSitDown(u, targetSeatId, seat)
	}

	resp.Room = room.ToCMSG()
}

func onReqRoomLookerInfo(u *user, req *cmsg.CReqRoomLookerInfo) {
	resp := &cmsg.SRespRoomLookerInfo{}
	defer u.SendMsg(resp)

	if !u.isInRoom() {
		resp.ErrCode = cmsg.SRespRoomLookerInfo_ErrNotInRoom
		return
	}
	room := u.room

	for _, v := range room.lookers {
		resp.Lookers = append(resp.Lookers, v.userBrief)
	}
}

func onReqRoomLeave(u *user, req *cmsg.CReqRoomLeave) {
	resp := &cmsg.SRespRoomLeave{}

	if !u.isInRoom() {
		resp.ErrCode = cmsg.SRespRoomLeave_ErrNotInRoom
		u.SendMsg(resp)
		return
	}

	if u.game != nil {
		u.game.reqQuit(u.userid, func(err error) {
			resp := &cmsg.SRespRoomLeave{}
			defer u.SendMsg(resp)
			if err != nil {
				resp.ErrCode = cmsg.SRespRoomLeave_ErrUnknown
				return
			}
			u.clearGame()
			u.room.quit(u, 0)
		})
	} else {
		u.room.quit(u, 0)
		u.SendMsg(resp)
		return
	}
}

func onReqRoomKick(u *user, req *cmsg.CReqRoomKick) {
	resp := &cmsg.SRespRoomKick{}
	defer u.SendMsg(resp)

	targetUserId := req.TargetUserId
	kickAllLooker := req.KickAllLooker

	if !u.isInRoom() {
		resp.ErrCode = cmsg.SRespRoomKick_ErrNotInRoom
		return
	}

	room := u.room

	if !room.isOwner(u.userid) {
		resp.ErrCode = cmsg.SRespRoomKick_ErrNoAuthority
		return
	}

	if room.isInGame() {
		resp.ErrCode = cmsg.SRespRoomKick_ErrRoomInGame
		return
	}

	if kickAllLooker {
		for _, user := range room.lookers {
			room.quit(user, u.userid)

			user.SendMsg(&cmsg.SNoticeRoomKick{
				KickerUserId: u.userid,
				RoomId:       room.roomId,
			})
		}
	} else {
		user, exist := room.findUser(targetUserId)
		if !exist {
			resp.ErrCode = cmsg.SRespRoomKick_ErrTargetNotInRoom
			return
		}

		room.quit(user, u.userid)

		user.SendMsg(&cmsg.SNoticeRoomKick{
			KickerUserId: u.userid,
			RoomId:       room.roomId,
		})
	}
}

func onReqRoomReady(u *user, req *cmsg.CReqRoomReady) {
	ready := req.Ready

	resp := &cmsg.SRespRoomReady{
		Ready: ready,
	}
	defer u.SendMsg(resp)

	if !u.isInRoom() {
		resp.ErrCode = cmsg.SRespRoomReady_ErrNotInRoom
		return
	}

	if u.isInGame() {
		resp.ErrCode = cmsg.SRespRoomReady_ErrRoomInGame
		return
	}

	room := u.room

	room.ready(u.userid, ready)

	if room.isAllReady() {
		room.createGame(func(game *Game, err error) {
			if err != nil {
				logrus.WithField("userId", u.userid).WithError(err).Error("room.createGame")
				return
			}
			room.clearAllReady()
		}, gameconf.GameStartTyp_GSTypeNormal)
	}
}

func onReqRoomStartGame(sender appframe.Session, req *cmsg.CReqRoomStartGame) {
	u, ok := userMgr.findUserBySessionID(sender.ID())
	if !ok {
		return
	}

	resp := &cmsg.SRespRoomStartGame{}

	if !u.isInRoom() {
		resp.ErrCode = cmsg.SRespRoomStartGame_ErrNotInRoom
		sender.SendMsg(resp)
		return
	}

	room := u.room

	if !room.isOwner(u.userid) {
		resp.ErrCode = cmsg.SRespRoomStartGame_ErrNoAuthority
		sender.SendMsg(resp)
	}

	if room.isInGame() {
		resp.ErrCode = cmsg.SRespRoomStartGame_ErrRoomInGame
		sender.SendMsg(resp)
		return
	}

	if !room.isAllReady() {
		resp.ErrCode = cmsg.SRespRoomStartGame_ErrNotAllReady
		sender.SendMsg(resp)
		return
	}

	room.createGame(func(game *Game, err error) {
		if err != nil {
			logrus.WithField("userId", u.userid).WithError(err).Error("room.createGame")
			resp.ErrCode = cmsg.SRespRoomStartGame_ErrCreateGameFail
			sender.SendMsg(resp)
			return
		}

		room.clearAllReady()

		sender.SendMsg(resp)
	}, gameconf.GameStartTyp_GSTypeNormal)
}
