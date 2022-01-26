package lobby

import (
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"math"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
	"time"
)

type Room struct {
	roomId  uint32
	roomNO  uint32
	voiceId string

	setting *gamedef.RoomSetting

	Owner   uint64
	Seats   []*RoomSeat
	lookers map[uint64]*user //旁观者

	game *Game

	waitReadyTimer *time.Timer
}

func (p *Room) noticeRoomEvent(event cmsg.SNoticeRoomEvent_Event, actionUser *user) {
	AppInstance.Post(func() {
		p.notifyMessage(&cmsg.SNoticeRoomEvent{
			Event:      event,
			ActionUser: actionUser.userBrief,
		}, nil)
	})
}
func (p *Room) noticeMatchGameStart() {
	notice := &cmsg.SNoticeMatchResult{
		Model:      int32(p.setting.GameMode),
		NoticeType: cmsg.SNoticeMatchResult_MatchStartGame,
		Room:       p.ToCMSG(),
	}
	for _, u := range p.Seats {
		if u.User != nil {
			u.User.SendMsg(notice)
		}
	}
}

func (p *Room) Enter(u *user, sitDown bool) (int32, bool) {
	if _, exist := p.findUser(u.userid); exist {
		return 0, false
	}

	p.noticeRoomEvent(cmsg.SNoticeRoomEvent_Enter, u)

	if sitDown {
		return p.sitDown(u)
	}

	p.sitLooker(u)
	return -1, true
}

func (p *Room) sitDown(u *user) (int32, bool) {
	for i, v := range p.Seats {
		if !v.isEmpty() {
			continue
		}

		seatId := int32(i)
		p.sitDownSeat(u, seatId, v)
		return seatId, true
	}
	return 0, false
}

func (p *Room) sitDownSeat(u *user, seatId int32, seat *RoomSeat) {
	seat.setUser(u)
	u.setRoom(p)
	p.notifyChange(cmsg.SNotifyRoomChange_SitDown, seatId, u.userid, 0, func(st *RoomSeat) bool {
		return st.User.userid == u.userid
	})

	if p.isSeatFull() {
		p.noticeRoomEvent(cmsg.SNoticeRoomEvent_SitFull, u)
	}

	if p.Owner == 0 {
		p.setOwner(u.userid)
		p.notifyChange(cmsg.SNotifyRoomChange_OwnerChange, seatId, u.userid, 0, func(st *RoomSeat) bool {
			return st.User.userid == u.userid
		})
	}
}

func (p *Room) sitLooker(u *user) bool {
	_, _, exist := p.findSeatByUserId(u.userid)
	if exist {
		return false
	}

	p.addLookerAndSync(u)
	u.setRoom(p)
	return true
}

func (p *Room) delLookerAndSync(userId uint64) {
	delete(p.lookers, userId)
	p.syncLookerCount()
}

func (p *Room) addLookerAndSync(u *user) {
	p.lookers[u.userid] = u
	p.syncLookerCount()
}

func (p *Room) ready(userId uint64, ready bool) bool {
	seatId, seat, exist := p.findSeatByUserId(userId)
	if !exist {
		return false
	}
	seat.setReady(ready)
	action := cmsg.SNotifyRoomChange_Ready
	if !ready {
		action = cmsg.SNotifyRoomChange_ReadyCancel
	}
	p.notifyChange(action, seatId, userId, 0, nil)
	return true
}

func (p *Room) isOpenReadyTimer() bool {
	if !p.isSeatFull() {
		return false
	}

	if p.isAllReady() {
		return false
	}

	var readyCount int32
	for _, v := range gameConfig.Global.NoticeReadyCondition {
		if v.K == int32(p.setting.MaxPlayer) {
			readyCount = v.V
			break
		}
	}

	return p.getReadyCount() >= readyCount
}

func (p *Room) checkReadyTimer() {
	open := p.isOpenReadyTimer()
	if open {
		if p.waitReadyTimer != nil {
			return
		}
		sec := gameConfig.Global.NoticeReadySec
		p.waitReadyTimer = time.AfterFunc(time.Duration(sec)*time.Second, func() {
			p.waitReadyTimer = nil
			AppInstance.Post(p.sitUpToLookerIfNotReady)
		})
		p.notifyMessage(&cmsg.SNoticeRoomTimer{
			NoticeType: cmsg.SNoticeRoomTimer_ReadyTimerStart,
			EndTime:    time.Now().Unix() + int64(sec),
		}, nil)
	} else {
		if p.waitReadyTimer == nil {
			return
		}

		p.StopReadyTimer()
	}
}

func (p *Room) StopReadyTimer() {
	if p.waitReadyTimer == nil {
		return
	}
	p.waitReadyTimer.Stop()
	p.waitReadyTimer = nil
	p.notifyMessage(&cmsg.SNoticeRoomTimer{
		NoticeType: cmsg.SNoticeRoomTimer_ReadyTimerCancel,
		EndTime:    0,
	}, nil)
}

//未准备的，变成旁观者
func (p *Room) sitUpToLookerIfNotReady() {
	if p.isInGame() {
		return
	}

	if p.isAllReady() {
		return
	}

	for i, seat := range p.Seats {
		if seat.isEmpty() {
			continue
		}
		if seat.isReady() {
			continue
		}

		p.SitUpToLooker(int32(i), seat)
	}
}

func (p *Room) SitUpToLooker(seatId int32, seat *RoomSeat) {
	u := seat.User
	seat.clearUser()

	p.sitLooker(u)
	u.SendMsg(&cmsg.SNoticeRoomChangeSeat{
		SeatId:     -1,
		LookerType: gameconf.LookerTyp_LTBlind,
	})
	p.notifyChange(cmsg.SNotifyRoomChange_SitUpToLooker, seatId, u.userid, 0, nil)

	p.changeOwnerIfOwnerLeave(u.userid, seatId)
}

func (p *Room) LookerSitDown(u *user, seatId int32, seat *RoomSeat) {
	p.delLookerAndSync(u.userid)
	p.sitDownSeat(u, seatId, seat)
}

func (p *Room) isUserReady(userId uint64) bool {
	_, seat, exist := p.findSeatByUserId(userId)
	if !exist {
		return false
	}
	return seat.isReady()
}

func (p *Room) changeSeatCount(count uint32) bool {
	for i, v := range p.Seats {
		seatId := uint32(i)
		if seatId >= count && !v.isEmpty() {
			return false
		}
	}

	seats := make([]*RoomSeat, count)
	for i := 0; i < int(count); i++ {
		if i < len(p.Seats) {
			seats[i] = p.Seats[i]
		} else {
			seats[i] = new(RoomSeat)
		}
	}
	p.Seats = seats
	return true
}

func (p *Room) updateSetting(cfg *gamedef.RoomSetting, userId uint64) {
	p.setting = cfg
	//更改设置后将座位准备状态置空
	for i := 0; i < len(p.Seats); i++ {
		p.Seats[i].Ready = false
	}
	seatId, _, _ := p.findSeatByUserId(userId)
	p.notifyChange(cmsg.SNotifyRoomChange_SettingChange, seatId, userId, 0, nil)
}

func (p *Room) isOwner(userId uint64) bool {
	return p.Owner == userId
}

func (p *Room) isAllReady() bool {
	return p.getReadyCount() == int32(len(p.Seats))
}

func (p *Room) getReadyCount() int32 {
	var count int32
	for _, v := range p.Seats {
		if v.isEmpty() {
			continue
		}

		if v.isReady() {
			count++
		}
	}
	return count
}

func (p *Room) clearAllReady() {
	p.setAllReady(false)
}

func (p *Room) setAllReady(ready bool) {
	for _, v := range p.Seats {
		if v.isEmpty() {
			continue
		}
		v.setReady(ready)
	}
}

func (p *Room) getUserCount() int32 {
	var count int32
	for _, v := range p.Seats {
		if !v.isEmpty() {
			count++
		}
	}
	return count
}

func (p *Room) SeatCount() int {
	return len(p.Seats)
}

// 获取房间用户ID
func (p *Room) getUserIds() []uint64 {
	seats := p.Seats
	res := make([]uint64, 0, len(seats))
	for _, v := range seats {
		if !v.isEmpty() {
			res = append(res, v.User.userid)
		}
	}
	return res
}

func (p *Room) isLooker(userId uint64) bool {
	_, exist := p.lookers[userId]
	return exist
}

func (p *Room) lookerQuit(u *user) {
	p.delLookerAndSync(u.userid)
	u.clearRoom()
}

func (p *Room) syncLookerCount() {
	p.notifyMessage(&cmsg.SSyncLookerCount{Count: int32(len(p.lookers))}, nil)
}

func (p *Room) quit(u *user, kickerUserId uint64) bool {
	defer p.noticeRoomEvent(cmsg.SNoticeRoomEvent_Leave, u)

	if p.isLooker(u.userid) {
		p.lookerQuit(u)
		p.delRoomIfEmpty()
		return true
	}

	seatId, _, _ := p.findSeatByUserId(u.userid)

	ok := p.removeUser(u, kickerUserId)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"roomId": p.roomId,
		}).Error("room remove user")
	}

	if del := p.delRoomIfEmpty(); del {
		return ok
	}

	p.changeOwnerIfOwnerLeave(u.userid, seatId)
	return ok
}

func (p *Room) changeOwnerIfOwnerLeave(actionUserId uint64, actionSeatId int32) {
	if p.Owner != actionUserId {
		return
	}

	var owner uint64 = 0
	if tu, _, exist := p.getFirstJoinUser(); exist {
		owner = tu.User.userid
	}

	p.setOwner(owner)
	p.notifyChange(cmsg.SNotifyRoomChange_OwnerChange, actionSeatId, actionUserId, 0, nil)
}

func (p *Room) getOneUser() (*RoomSeat, int32, bool) {
	for i, seat := range p.Seats {
		if seat.isEmpty() {
			continue
		}

		return seat, int32(i), true
	}
	return nil, 0, false
}

func (p *Room) getFirstJoinUser() (*RoomSeat, int32, bool) {
	var minJoinTime int64 = math.MaxInt64

	for _, seat := range p.Seats {
		if seat.isEmpty() {
			continue
		}

		if seat.JoinTime < minJoinTime {
			minJoinTime = seat.JoinTime
		}
	}

	for i, seat := range p.Seats {
		if seat.isEmpty() {
			continue
		}

		if seat.JoinTime == minJoinTime {
			return seat, int32(i), true
		}
	}

	return nil, 0, false
}

func (p *Room) getUsers() (users []*user) {
	for _, seat := range p.Seats {
		if seat.isEmpty() {
			continue
		}
		users = append(users, seat.User)
	}
	return
}

func (p *Room) delRoomIfEmpty() bool {
	if !p.isEmpty() {
		return false
	}

	roomMgr.delRoomById(p.roomId)
	return true
}

func (p *Room) isEmpty() bool {
	return p.isAllSeatEmpty() && len(p.lookers) == 0
}

func (p *Room) isAllSeatEmpty() bool {
	for _, v := range p.Seats {
		if !v.isEmpty() {
			return false
		}
	}
	return true
}

func (p *Room) isAllowEnter() bool {
	return p.setting.AllowEnter
}

//有空座位而且允许进
func (p *Room) canSeatNewPlayer() bool {
	if !p.isAllowEnter() {
		return false
	}
	if p.isSeatFull() {
		return false
	}
	if p.isInGame() {
		return false
	}
	return true
}

func (p *Room) isInGame() bool {
	return p.game != nil
}

func (p *Room) setGame(game *Game) {
	p.game = game
	for _, v := range p.Seats {
		if v.isEmpty() {
			continue
		}
		v.User.bindGame(game)
	}

	for _, v := range p.lookers {
		v.bindGame(game)
	}
}

func (p *Room) clearGame() {
	gsTyp := p.game.gameStartType
	p.game = nil
	for _, v := range p.Seats {
		if v.isEmpty() {
			continue
		}
		v.User.clearGame()
		v.User.setUserGameStatus(gameconf.UserGameStatusTyp_UGSTFree)
		if gsTyp == gameconf.GameStartTyp_GSTypeMatch {
			v.User.quitRoomIfMatch()
		} else {
			v.User.quitRoomIfOffline()
		}
	}

	for _, v := range p.lookers {
		v.clearGame()
		v.quitRoomIfOffline()
	}
}

func (p *Room) isSeatFull() bool {
	for _, v := range p.Seats {
		if v.isEmpty() {
			return false
		}
	}
	return true
}

func (p *Room) findSeatBySeatId(seatId int32) (*RoomSeat, bool) {
	for i, seat := range p.Seats {
		if seatId == int32(i) {
			return seat, true
		}
	}
	return nil, false
}

func (p *Room) findSeatId(userId uint64) (int32, bool) {
	if p.isLooker(userId) {
		return -1, true
	}

	seatId, _, ok := p.findSeatByUserId(userId)
	if ok {
		return seatId, true
	}

	return 0, false
}

func (p *Room) findSeatByUserId(userId uint64) (int32, *RoomSeat, bool) {
	for i, seat := range p.Seats {
		u, exist := seat.getUser()
		if !exist || userId != u.userid {
			continue
		}
		return int32(i), seat, true
	}
	return 0, nil, false
}

func (p *Room) findUser(userId uint64) (*user, bool) {
	_, seat, ok := p.findSeatByUserId(userId)
	if ok {
		return seat.User, true
	}
	user, ok := p.lookers[userId]
	if ok {
		return user, true
	}
	return nil, false
}

func (p *Room) removeUser(u *user, kickerUserId uint64) bool {
	var remove bool
	for i, seat := range p.Seats {
		roomUser, exist := seat.getUser()
		if !exist || u.userid != roomUser.userid {
			continue
		}
		seat.clearUser()
		remove = true
		u.clearRoom()

		action := cmsg.SNotifyRoomChange_Leave
		if kickerUserId != 0 {
			action = cmsg.SNotifyRoomChange_LeaveByKick
		}
		p.notifyChange(action, int32(i), u.userid, kickerUserId, nil)
		break
	}
	return remove
}

func (p *Room) notifyMessage(msg proto.Message, filter func(seat *RoomSeat) bool) {
	for _, v := range p.Seats {
		if v.isEmpty() {
			continue
		}
		if filter != nil && filter(v) {
			continue
		}
		v.User.SendMsg(msg)
	}

	for _, u := range p.lookers {
		u.SendMsg(msg)
	}
}

func (p *Room) clearUser() {
	for _, seat := range p.Seats {
		if u, ok := seat.getUser(); ok {
			seat.clearUser()
			u.clearRoom()
		}
	}

	for _, u := range p.lookers {
		u.clearRoom()
	}
}

func (p *Room) setOwner(owner uint64) {
	p.Owner = owner
}

func (p *Room) ToCMSG() *gamedef.Room {
	return &gamedef.Room{
		OwnerUserId: p.Owner,
		Setting:     p.setting,
		RoomId:      p.roomId,
		RoomNO:      p.roomNO,
		VoiceId:     p.voiceId,
		LookerCount: int32(len(p.lookers)),
		Seats:       p.seatInfos(),
	}
}

func (p *Room) onUserDisconnect(u *user) {
	if u.game != nil {
		if p.isLooker(u.userid) {
			u.game.reqQuit(u.userid, func(err error) {
				if err != nil {
					logrus.WithError(err).WithField("userId", u.userid).Error("u.game.reqQuit")
					return
				}
				u.clearGame()
				u.room.quit(u, 0)
			})
		} else {
			u.game.node.SendMsg(&smsg.UserDisconnect{
				Userid: u.userid,
			})
		}
		return
	}
	if u.GameStatus == gameconf.UserGameStatusTyp_UGSTReadyRoom {
		return
	}

	u.room.quit(u, 0)
}

func (p *Room) seatInfos() []*gamedef.Room_Seat {
	res := make([]*gamedef.Room_Seat, 0, len(p.Seats))
	for i, v := range p.Seats {
		seatId := int32(i)
		if !v.isEmpty() {
			u := v.User
			res = append(res, &gamedef.Room_Seat{
				SeatId:    seatId,
				UserBrief: u.userBrief,
				Ready:     v.Ready,
			})
		} else {
			res = append(res, &gamedef.Room_Seat{
				SeatId:    seatId,
				UserBrief: nil,
				Ready:     v.Ready,
			})
		}
	}
	return res
}

func (p *Room) notifyChange(action cmsg.SNotifyRoomChange_Action, actionSeatId int32, actionUserId uint64, kicker uint64, filter func(seat *RoomSeat) bool) {
	p.notifyMessage(&cmsg.SNotifyRoomChange{
		Action:       action,
		ActionUserId: actionUserId,
		ActionSeatId: actionSeatId,
		KickerUserId: kicker,
		Room:         p.ToCMSG(),
	}, filter)

	AppInstance.Post(p.checkReadyTimer)
}

func (p *Room) createGame(callback func(*Game, error), gsType gameconf.GameStartTyp) {
	cbk := func(g *Game, err error) {
		callback(g, err)

		if err != nil {
			p.clearGame()
			return
		}
	}

	var aiSvrId uint32
	g, err := gameMgr.createGame(p.setting.GameMode, p.roomId, aiSvrId, len(p.Seats), gsType, p.getUserIds())
	if err != nil {
		cbk(nil, err)
		return
	}

	p.setGame(g)

	users := make([]*smsg.LoGaReqNewGame_UserInfo, 0, len(p.Seats))
	for i, seat := range p.Seats {
		if seat.isEmpty() {
			continue
		}
		seatId := int32(i)
		u := seat.User
		user := &smsg.LoGaReqNewGame_UserInfo{
			UserId:    u.userid,
			GateId:    u.session.SvrID,
			Session:   u.session.ID,
			IsRobot:   u.isRobot,
			UserBrief: u.userBrief,
			SeatId:    seatId,
		}
		users = append(users, user)
	}

	lookers := make([]*smsg.LoGaReqNewGame_UserInfo, 0, len(p.lookers))
	for _, u := range p.lookers {
		looker := &smsg.LoGaReqNewGame_UserInfo{
			UserId:    u.userid,
			GateId:    u.session.SvrID,
			Session:   u.session.ID,
			IsRobot:   u.isRobot,
			UserBrief: u.userBrief,
			SeatId:    -1,
		}
		lookers = append(lookers, looker)
	}

	req := &smsg.LoGaReqNewGame{
		GameUUID:    g.gameId,
		Users:       users,
		Lookers:     lookers,
		AiSvrID:     aiSvrId,
		RoomId:      p.roomId,
		RoomNO:      p.roomNO,
		VoiceId:     p.voiceId,
		RoomSetting: p.setting,
	}

	g.node.ReqSugar(req, func(resp *smsg.LoGaRespNewPVPGame, err error) {
		log := logrus.WithField(" node : ", g.node.ID())
		if err != nil {
			gameMgr.deleteGame(g.gameId)
			log.WithError(err).Error("New Game failed")
			cbk(nil, err)
			return
		}

		if resp.ErrCode != 0 {
			gameMgr.deleteGame(g.gameId)
			log.WithError(err).Error("New Game failed ", resp.ErrCode)
			cbk(nil, err)
			return
		}

		log.WithField("GameID", g.gameId).Debug("New Game succ")
		cbk(g, nil)
	}, time.Second*20)
}

func (p *Room) lookGame(u *user) {
	looker := &smsg.LoGaReqLookGame_UserInfo{
		UserId:    u.userid,
		GateId:    u.session.SvrID,
		Session:   u.session.ID,
		IsRobot:   u.isRobot,
		UserBrief: u.userBrief,
		SeatId:    -1,
	}

	req := &smsg.LoGaReqLookGame{
		GameUUID: p.game.gameId,
		Looker:   looker,
		RoomId:   p.roomId,
		RoomNO:   p.roomNO,
	}

	u.bindGame(p.game)

	p.game.node.ReqSugar(req, func(resp *smsg.LoGaRespLookGame, err error) {
		log := logrus.WithField(" node : ", p.game.node.ID())
		if err != nil {
			log.WithError(err).Error("Look Game failed")
			return
		}
		if resp.ErrCode != 0 {
			log.WithError(err).Error("Look Game failed ", resp.ErrCode)
			return
		}
		log.WithField("GameID", p.game.gameId).Debug("Look Game succ")
	}, time.Second*20)
}

func (p *Room) joinPlayer(userId uint64) {
	u, ok := userMgr.findUser(userId)
	if !ok {
		logrus.Warnf("find user failed,user id %d", userId)
		return
	}
	_, ok = p.Enter(u, true)
	if !ok {
		logrus.Warnf("player enter room failed, user id %d", userId)
		return
	}
}

func (p *Room) checkReadyStartGame() {
	if p.isAllReady() {
		p.createGame(func(game *Game, err error) {
			if err != nil {
				logrus.WithField("roomNo", p.roomNO).WithError(err).Error("room.createGame")
				return
			}
			p.clearAllReady()
		}, gameconf.GameStartTyp_GSTypeMatch)
	}
}

type RoomSeat struct {
	User     *user
	Ready    bool
	JoinTime int64 //加入时间
}

func (p *RoomSeat) isEmpty() bool {
	return p.User == nil
}

func (p *RoomSeat) isReady() bool {
	return p.Ready
}

func (p *RoomSeat) setReady(ready bool) {
	p.Ready = ready
}

func (p *RoomSeat) setUser(u *user) {
	p.User = u
	if u != nil {
		p.JoinTime = time.Now().Unix()
	}
}

func (p *RoomSeat) clearUser() {
	p.User = nil
	p.Ready = false
}

func (p *RoomSeat) getUser() (*user, bool) {
	if p.User == nil {
		return nil, false
	}
	return p.User, true
}
