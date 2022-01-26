package lobby

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type MatchQueue struct {
	UserID       uint64 //单人或房主userID
	QueueEndTime int64
	Users        []*UserMatchParam //多人时，是数组
	CurQueueTyp  MatchQueueTyp

	ModeID    int32
	JoinTime  int64
	UserNum   int
	StartTime int64 //开始匹配时间

	TeamID uint32

	MatchScore int32
}

func (p *MatchQueue) getUserNum() int32 {
	return int32(p.UserNum)
}

func (p *MatchQueue) isTeam() bool {
	return p.TeamID > 0 && p.UserNum > 1
}

func (p *MatchQueue) SetMatch(match bool, matchModeID int32, action Action) bool {
	if p.TeamID == 0 {
		us, exist := userMgr.findUser(p.UserID)
		if !exist {
			return false
		}
		us.SetMatch(match, matchModeID)
		return true
	}
	return true
}

func (p *MatchQueue) SetGameStatus(gameStatus gameconf.UserGameStatusTyp) bool {
	if p.TeamID == 0 {
		us, exist := userMgr.findUser(p.UserID)
		if !exist {
			return false
		}
		us.setUserGameStatus(gameStatus)
		return true
	}
	return true
}

func (p *MatchQueue) getUsers() []*UserMatchParam {
	return p.Users
}

func (p *MatchQueue) notifyMessage(msg proto.Message) error {
	if p.TeamID == 0 {
		user, exist := userMgr.findUser(p.UserID)
		if exist {
			user.SendMsg(msg)
		}
		return nil
	}
	return nil
}

type MatchRoom struct {
	Mqs   map[uint64]*MatchQueue
	Count int32
	Star  int32
}

func (p *MatchRoom) init() {
	p.Mqs = make(map[uint64]*MatchQueue)
	p.Count = 0
}

func (p *MatchRoom) getCount() int32 {
	return p.Count
}

func (p *MatchRoom) add(queue *MatchQueue) {
	p.Mqs[queue.UserID] = queue
	p.Count += queue.getUserNum()
}

func (p *MatchRoom) setMatch(isMatch bool, modeID int32, action Action) {
	for _, v := range p.Mqs {
		v.SetMatch(isMatch, modeID, action)
	}
}

func (p *MatchRoom) setGameStatus(gameStatus gameconf.UserGameStatusTyp) {
	for _, v := range p.Mqs {
		v.SetGameStatus(gameStatus)
	}
}

func (p *MatchRoom) getUserCount() int32 {
	return p.Count
}
