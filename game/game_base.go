package game

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"time"
)

type GameModeInterface interface {
	//InitFSM(base *GameBase) *fsm.FSM
	OnEnterState(state string)
	//LeaveState(state string)
	CheckIsOver() bool
}

type GameBase struct {
	core.LogicCore
	triggerMgr *TriggerManager

	gameStartTime time.Time
	gameEndTime   time.Time

	setting *gamedef.RoomSetting
	config  *conf.GameConfig

	players   []*Player
	seatCount int32
	lookers   map[uint64]*Player

	gameUUID   string
	gameMode   gameconf.GameModeTyp
	roomId     uint32
	roomNO     uint32
	roundCount int32

	Phase        gamedef.GamePhase
	PhaseEndTime int64

	opSeatId  int32 //当前正在操作的座位号,语音场描述阶段用
	opEndTime int64 //当前操作的截止时间,语音场描述阶段用

	curPlayer int32

	over          bool
	clearCallback func()
	startCallback func()

	modelImp GameModeInterface

	board *Board
}

func (g *GameBase) Stop() {
	g.FSMEvent("stop")
}

func (g *GameBase) Clear() {
	g.clearCallback()
}

func newGameBase(worker core.Worker, setting *gamedef.RoomSetting, gameUUID string, roomId, roomNO uint32, config *conf.GameConfig) *GameBase {
	now := time.Now()

	g := &GameBase{
		triggerMgr:    nil,
		gameStartTime: now,
		gameEndTime:   now,
		setting:       setting,
		players:       make([]*Player, setting.MaxPlayer),
		lookers:       make(map[uint64]*Player, 0),
		gameUUID:      gameUUID,
		gameMode:      setting.GameMode,
		roomId:        roomId,
		roomNO:        roomNO,
		roundCount:    0,
		Phase:         0,
		over:          false,
		config:        config,
		board:         newBoard(),
		seatCount:     int32(setting.MaxPlayer),
	}

	g.LogicCore.Init(worker)

	g.triggerMgr = newTriggerManager(g)
	g.triggerMgr.init()
	return g
}

func (g *GameBase) GetTriggerMgr() *TriggerManager {
	return g.triggerMgr
}

func (g *GameBase) TryEndGame() bool {
	if g.modelImp.CheckIsOver() {
		g.Stop()
		return true
	}

	return false
}

func (g *GameBase) OnUserJoin(u *User, seatId int32) *Player {
	player := newPlayer(u, g, seatId)
	g.players[seatId] = player
	return player
}

func (g *GameBase) onLookerJoin(u *User) *Player {
	looker := newPlayer(u, g, -1)
	g.lookers[u.userId] = looker
	return looker
}

func (g *GameBase) Send2All(msg proto.Message) {
	g.notifyMessage(msg, nil)
	g.sendToLookers(msg)
}

func (g *GameBase) SendSeatMsg(toMsg func(seatId int32) proto.Message) {
	for _, v := range g.players {
		if v == nil {
			continue
		}

		msg := toMsg(v.seatId)
		v.SendMsg(msg)
	}
}

func (g *GameBase) notifyMessage(msg proto.Message, filter func(*Player) bool) {
	for _, v := range g.players {
		if v == nil {
			continue
		}

		if filter != nil && filter(v) {
			continue
		}
		v.SendMsg(msg)
	}
}

func (p *GameBase) alivePlayers() []*Player {
	players := make([]*Player, 0)
	for _, v := range p.players {
		if v == nil {
			continue
		}
		if v.dead {
			continue
		}
		players = append(players, v)
	}
	return players
}

func (p *GameBase) aliveSeatIds() []int32 {
	seatIds := make([]int32, 0)
	for _, v := range p.players {
		if v == nil {
			continue
		}
		if v.dead {
			continue
		}
		seatIds = append(seatIds, v.seatId)
	}
	return seatIds
}

func (g *GameBase) getPlayerBySeat(seatId int32) (*Player, bool) {
	for _, v := range g.players {
		if v == nil {
			continue
		}
		if v.seatId == seatId {
			return v, true
		}
	}

	return nil, false
}

func (g *GameBase) isInPlayer(p *Player, players []*Player) bool {
	for _, v := range players {
		if v == p {
			return true
		}
	}
	return false
}

func (g *GameBase) findPlayerByUserId(userId uint64) (*Player, bool) {
	for _, v := range g.players {
		if v.user.userId == userId {
			return v, true
		}
	}
	return nil, false
}

func (g *GameBase) SetOver() {
	g.over = true
}

func (g *GameBase) IsOver() bool {
	return g.over
}

//强制结束
func (g *GameBase) ForceOver() {
	g.SetOver()
	g.Clear()
}

func (g *GameBase) SetAllClientReady(v bool) {
	for _, p := range g.players {
		if p == nil {
			continue
		}
		p.setClientReady(v)
	}
	for _, p := range g.lookers {
		if p == nil {
			continue
		}
		p.setClientReady(v)
	}
}

func (g *GameBase) isInAlivePlayer(p *Player) bool {
	alive := g.alivePlayers()
	for i := 0; i < len(alive); i++ {
		if p == alive[i] {
			return true
		}
	}
	return false
}

func (g *GameBase) syncPlayerState(player *Player) {
	SyncPlayerState(player, cmsg.SSyncPlayerState_State)
}

func (g *GameBase) GetAliveCount(typ gameconf.RoleTyp) int32 {
	var count int32
	for _, v := range g.players {
		if v == nil {
			continue
		}
		if v.IsDead() {
			continue
		}
		if v.roleType == typ {
			count++
		}
	}
	return count
}

func (g *GameBase) GetAlivePlayers() []*Player {
	var players []*Player
	for _, v := range g.players {
		if v == nil {
			continue
		}
		if v.IsDead() {
			continue
		}
		players = append(players, v)
	}
	return players
}

func (g *GameBase) GetCurrentPlayer() *Player {
	return g.players[g.curPlayer]
}

func (g *GameBase) seatInfos(players ...*Player) []*gamedef.GameSeat {
	var seats []*gamedef.GameSeat

	if len(players) == 0 {
		for _, v := range g.players {
			if v == nil {
				continue
			}
			seats = append(seats, v.toSeatInfo(true))
		}
	} else {
		for _, v := range players {
			seats = append(seats, v.toSeatInfo(true))
		}
	}

	return seats
}

func (g *GameBase) setPhaseAndSync(phase gamedef.GamePhase) int32 {
	g.Phase = phase
	var sec int32 = 20
	switch phase {
	case gamedef.GamePhase_PhaseMain:
		sec = g.config.OperationTime
	default:

	}

	g.PhaseEndTime = time.Now().Unix() + int64(sec)
	g.Send2All(&cmsg.SNoticeEnterPhase{
		Phase:        g.Phase,
		RoundCount:   g.roundCount,
		SeatId:       g.GetCurrentPlayer().seatId,
		PhaseEndTime: g.PhaseEndTime,
	})
	return sec
}

func (g *GameBase) onPlayerDisconnect(p *Player) {
	p.disconnect()
	g.syncPlayerState(p)
	g.resetWaitTimeIfOpPlayerLeave(p)
}

func (g *GameBase) resetWaitTimeIfOpPlayerLeave(p *Player) {
	//if !g.waitMgr.isWait(p, &cmsg.CReqDescribe{}) {
	//	return
	//}
	//
	//lastOp, exist := g.waitMgr.getPlayerLastOp(p)
	//if !exist {
	//	return
	//}
	//
	//opMsg, ok := lastOp.msg.(*cmsg.SNoticeOp)
	//if !ok {
	//	return
	//}
	//
	//opWaitTime := time.Duration(g.config.OpSecWhenUserLeave) * time.Second
	//
	//opEndTime, ok := g.waitMgr.resetIfWaitTimeLong(opWaitTime)
	//if !ok {
	//	return
	//}
	//
	//for _, v := range g.players {
	//	opMsg := &cmsg.SNoticeOp{
	//		OpType:    opMsg.OpType,
	//		OpSeatIds: opMsg.OpSeatIds,
	//		Word:      v.word,
	//		OpEndTime: opEndTime,
	//	}
	//	v.SendMsg(opMsg)
	//}
	//
	//g.sendToLookers(&cmsg.SNoticeOp{
	//	OpType:    opMsg.OpType,
	//	OpSeatIds: opMsg.OpSeatIds,
	//	Word:      "",
	//	OpEndTime: opEndTime,
	//})
}

func (g *GameBase) onPlayerReconnect(p *Player, session appframe.Session) {
	p.reconnect(session)
	g.syncPlayerState(p)
}

func (g *GameBase) isLooker(userId uint64) bool {
	_, exist := g.lookers[userId]
	return exist
}

func (g *GameBase) onPlayerQuit(p *Player) {
	p.quit()
	g.syncPlayerState(p)
	g.resetWaitTimeIfOpPlayerLeave(p)
}

func (g *GameBase) onLookerQuit(p *Player) {
	delete(g.lookers, p.user.userId)
}

func (g *GameBase) resendLastOpMsg(p *Player) {
	lastOp, exist := g.LogicCore.GetPlayerLastOp(p)
	if !exist {
		return
	}

	p.SendMsg(lastOp.Msg)
}

func (g *GameBase) checkStateSync() {
	for _, v := range g.players {
		if v.user.session.ID().ID == 0 {
			v.connectState = gameconf.UserConnectState_USDisconnect
			g.syncPlayerState(v)
		}
	}
}

func (g *GameBase) prepareStart() {
	g.FSMEvent("prepare")
}

func (g *GameBase) CanActionPlayers() (players []*Player) {
	for _, v := range g.players {
		if g.canPlayerAction(v) {
			players = append(players, v)
		}
	}
	return
}

func (g *GameBase) canPlayerAction(player *Player) bool {
	if g.board.SeatCardCount(player.GetSeatID()) > 0 {
		return true
	}

	//if player.HandCard.Count()+player.HandCardPool.Count() > 0 && g.board.SeatEmptyCellCount(player.GetSeatID()) > 0 {
	//	return true
	//}

	if player.HandCard.Count()+player.HandCardPool.Count() > 0 {
		return true
	}

	return false
}

func (g *GameBase) noticeGameStart() {
	for _, v := range g.players {
		if v == nil {
			continue
		}
		v.SendMsg(&cmsg.SNoticeGameStart{})
	}

	g.sendToLookers(&cmsg.SNoticeGameStart{})
}

func (g *GameBase) sendToLookers(msg proto.Message) {
	for _, u := range g.lookers {
		u.SendMsg(msg)
	}
}

func (g *GameBase) GetNextAlivePlayer(curSeat int32) *Player {
	for seat := (curSeat + 1) % g.seatCount; seat != curSeat; seat = (seat + 1) % g.seatCount {
		if g.players[seat].IsDead() {
			continue
		}
		return g.players[seat]
	}

	return nil
}

func (g *GameBase) NoticeOpCommon(opPlayer *Player, opType cmsg.SNoticeOp_OpType, waitSec int64, msgAllow []proto.Message, waitCallback func(timeout bool)) {
	now := time.Now()
	opEndTime := now.Unix() + waitSec

	msg := &cmsg.SNoticeOp{
		OpType:    opType,
		OpSeatId:  opPlayer.GetSeatID(),
		TargetPos: nil,
		OpEndTime: opEndTime,
		LeftAP:    opPlayer.GetAP(),
	}
	g.Send2All(msg)
	g.StartWaiting([]core.Player{opPlayer}, waitSec, msg, msgAllow, waitCallback)
}

func (g *GameBase) phaseLeftSec() int64 {
	return g.PhaseEndTime - time.Now().Unix()
}
