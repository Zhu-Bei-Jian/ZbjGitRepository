package lobby

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"time"
)

type Mode1MatchManager struct {
	*MatchManager
	//队列
	UserID2MatchQueue map[uint64]*MatchQueue
	Tick              *time.Ticker
	ModeID            int32
	MaxPlayerCnt      int32

	// 准备桌子管理
	IdMaker uint32
}

func (mm *Mode1MatchManager) deleteMatchQueue(uid uint64) {
	delete(mm.UserID2MatchQueue, uid)
}
func (p *Mode1MatchManager) init(modeID int32, matchMgr *MatchManager) {
	p.MatchManager = matchMgr
	p.UserID2MatchQueue = make(map[uint64]*MatchQueue)
	p.ModeID = modeID
	p.MaxPlayerCnt = 2
}

func (p *Mode1MatchManager) joinMatch(userID uint64, score int32) (*MatchQueue, error) {
	//u, ok := userMgr.findUser(userID)
	//if !ok {
	//	return nil, errors.Errorf("join match find user failed,user id %v", userID)
	//}
	//
	//if u.GameStatus == gameconf.UserGameStatusTyp_UGSTReadyRoom {
	//	return nil, errors.Errorf("user %v in ready room", userID)
	//}

	mq, err := p.change2TypeQueue(userID, score, MatchQueueTyp_MQTWaitQueue)
	if err != nil {
		return nil, err
	}

	//if !p.MatchManager.Testing {
	//	mq.SetMatch(true, p.ModeID, Action_Match)
	//	mq.SetGameStatus(gameconf.UserGameStatusTyp_UGSTMatching)
	//}
	mq.StartTime = gameutil.GetCurrentTimestamp()
	return mq, nil
}

func (p *Mode1MatchManager) change2TypeQueue(userID uint64, score int32, queTyp MatchQueueTyp) (*MatchQueue, error) {
	_, exist := p.getMatchQueueByUserID(userID)
	if exist {
		return nil, errAlreadyInMatch
	}
	mq := &MatchQueue{
		UserID:       userID,
		QueueEndTime: p.getCurQueueWaitSeconds(p.ModeID, queTyp),
		JoinTime:     gameutil.GetCurrentTimestamp(),
		ModeID:       p.ModeID,
		MatchScore:   score,
		CurQueueTyp:  queTyp,
		UserNum:      1,
	}
	p.UserID2MatchQueue[userID] = mq
	return mq, nil
}

func (p *Mode1MatchManager) getMatchQueueByUserID(userID uint64) (*MatchQueue, bool) {
	mq, exist := p.UserID2MatchQueue[userID]
	return mq, exist
}

func (p *Mode1MatchManager) match() {
	now := time.Now()
	defer func() {
		elaspe := time.Since(now)
		if elaspe > time.Millisecond {
			fmt.Println("matchmode1 elaspe:", elaspe)
		}
	}()
	//匹配超时的删除
	//p.Limit.GC(p.ModeID)
	res := make([]*MatchRoom, 0)
	res = append(res, p.NoRuleMatch()...)

	if !p.MatchManager.Testing {
		p.createGame(res)
	} else {
		now := gameutil.GetCurrentTimestamp()
		for _, v := range res {
			str := ""
			for _, mq := range v.Mqs {
				str += fmt.Sprintf("matchTime:%v user id %v\n", now-mq.StartTime, mq.UserID)
			}
			fmt.Println(str)
		}
	}
}

//不开匹配规则的任意匹配
func (p *Mode1MatchManager) NoRuleMatch() []*MatchRoom {
	res := make([]*MatchRoom, 0)
	for _, currentQue := range p.UserID2MatchQueue {
		matchSuccess := false
		mt := &MatchRoom{}
		mt.init()
		mt.add(currentQue)
		for _, mqs := range p.UserID2MatchQueue {
			if currentQue.UserID == mqs.UserID {
				continue
			}
			if mt.getCount() < 2 {
				mt.add(mqs)
			}

			if mt.getCount() == 2 {
				matchSuccess = true
				goto MatchSuccess
			}
		}
	MatchSuccess:
		if matchSuccess {
			res = append(res, mt)
			for _, v := range mt.Mqs {
				p.deleteMatchQueue(v.UserID)
			}
		}
	}
	return res
}

func (p *Mode1MatchManager) getQualifyRankRange(rank uint32, rg uint32) (uint32, uint32) {
	var min, max uint32
	if rank > rg {
		min = rank - rg
	}
	max = rank + rg
	return min, max
}

func (p *Mode1MatchManager) matchingUserCount() int32 {
	return int32(len(p.UserID2MatchQueue))
}

func (p *Mode1MatchManager) createGame(mrs []*MatchRoom) {
	var room *Room
	var err error
	for _, mr := range mrs {
		p.noticeMatchResult(mr)
		room, err = p.createReadyRoom(mr)
		if err != nil {
			logrus.WithError(err).Error("create ready room error")
			return
		}
		for _, player := range mr.Mqs {
			room.joinPlayer(player.UserID)
			room.ready(player.UserID, true)
		}
		room.Owner = 0
		mr.setGameStatus(gameconf.UserGameStatusTyp_UGSTReadyRoom)
		mr.setMatch(false, p.ModeID, Action_EnterReadyRoom)
	}

	for _, mr := range mrs {
		mr.setGameStatus(gameconf.UserGameStatusTyp_UGSTGame)
		room.noticeMatchGameStart()
		room.checkReadyStartGame()
	}

}
func (p *Mode1MatchManager) noticeMatchResult(mr *MatchRoom) {
	var players []*gamedef.UserBrief
	for _, v := range mr.Mqs {
		user, ok := userMgr.findUser(v.UserID)
		if !ok {
			logrus.Warnf("find ueser failed,user id %d", v.UserID)
			continue
		}
		players = append(players, user.userBrief)
	}

	msg := &cmsg.SNoticeMatchResult{
		NoticeType: cmsg.SNoticeMatchResult_MatchSucess,
		Model:      p.ModeID,
		Players:    players,
	}

	for _, v := range mr.Mqs {
		user, ok := userMgr.findUser(v.UserID)
		if !ok {
			logrus.Warnf("find ueser failed,user id %d", v.UserID)
			continue
		}
		user.SendMsg(msg)
	}
}

func (p *Mode1MatchManager) createReadyRoom(mr *MatchRoom) (*Room, error) {
	setting := &gamedef.RoomSetting{
		GameMode:   gameconf.GameModeTyp(p.ModeID),
		RoomName:   "",
		MaxPlayer:  uint32(len(mr.Mqs)),
		AllowEnter: false,
	}

	var ownerId uint64
	for _, v := range mr.Mqs {
		ownerId = v.UserID
		room, err := roomMgr.newRoom(setting, ownerId)
		if err != nil {
			logrus.WithError(err).Errorf("create room failed,user id %d", ownerId)
			return nil, err
		}
		return room, err
	}
	return nil, errors.New("create room failed")
}

func (p *Mode1MatchManager) newID() uint32 {
	p.IdMaker++
	return p.IdMaker
}

func (p *Mode1MatchManager) quitMatch(userID uint64) (*MatchQueue, error) {
	mq, exist := p.getMatchQueueByUserID(userID)
	if !exist {
		return nil, errNotInMatch
	}
	mq.SetMatch(false, 0, Action_MatchCancel)
	p.deleteMatchQueue(mq.UserID)
	return mq, nil
}

func (p *Mode1MatchManager) run() {
	p.Tick = workerMgr.Ticker(time.Millisecond*400, p.match)
}
func (p *Mode1MatchManager) close() {
	if p.Tick != nil {
		p.Tick.Stop()
	}
}
