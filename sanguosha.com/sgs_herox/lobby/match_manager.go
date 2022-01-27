package lobby

import (
	"github.com/pkg/errors"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

var errAlreadyInMatch = errors.New("user already in match")
var errNotInMatch = errors.New("user not in match")

type MatchQueueTyp int32

const (
	MatchQueueTyp_MQTInvalid MatchQueueTyp = 0
	// 等待匹配队列
	MatchQueueTyp_MQTWaitQueue MatchQueueTyp = 1
	// 严格匹配队列
	MatchQueueTyp_MQTStrictQueue MatchQueueTyp = 2
)

type Matcher interface {
	run()
	close()
	match()
	quitMatch(userID uint64) (*MatchQueue, error)
	joinMatch(userID uint64, score int32) (*MatchQueue, error)
	matchingUserCount() int32 //还在匹配队列中的人数
}

type MatchManager struct {
	Mode2ModeMatchManager map[int32]Matcher
	Testing               bool
}

func (p *MatchManager) init() {
	p.Mode2ModeMatchManager = make(map[int32]Matcher)

	{
		modeId := int32(gameconf.GameModeTyp_MGTSpyText)
		matcher := new(Mode1MatchManager)
		matcher.init(modeId, p)
		p.Mode2ModeMatchManager[modeId] = matcher
	}

}

type UserMatchParam struct {
	WinRate        int32 //胜率
	ConsecutiveCnt int32 //连胜连败
	LoseC11Cnt     int32 //连败场次（身份场）
	score          int32
	userBrief      *gamedef.UserBrief
}

func (p *MatchManager) joinMatch(typ MatchQueueTyp, userID uint64, modeID int32, score int32) error {
	mm, err := p.getModeMatchManager(modeID)
	if err != nil {
		return err
	}

	_, err = mm.joinMatch(userID, score)
	if err != nil {
		return err
	}
	return nil
}

func (p *MatchManager) getModeMatchManager(modeID int32) (Matcher, error) {
	mm, exist := p.Mode2ModeMatchManager[modeID]
	if !exist {
		return nil, errors.New("mode not exist")
	}
	return mm, nil
}

func (p *MatchManager) run() {
	for _, v := range p.Mode2ModeMatchManager {
		v.run()
	}
}

func (p *MatchManager) close() {
	for _, v := range p.Mode2ModeMatchManager {
		v.close()
	}
}

func (p *MatchManager) quitMatch(us *user) (bool, error) {
	//modeID := 1
	mm, err := p.getModeMatchManager(1)
	if err != nil {
		return false, err
	}
	userID := us.userid
	_, err = mm.quitMatch(userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *MatchManager) getCurQueueWaitSeconds(modeID int32, typ MatchQueueTyp) int64 {
	//暂时就一个排位赛，简单处理
	var waitSecond int32
	switch typ {
	case MatchQueueTyp_MQTWaitQueue:
		waitSecond = 100
	case MatchQueueTyp_MQTStrictQueue:
		waitSecond = 100000000
	}
	queueEndTime := int64(waitSecond) + gameutil.GetCurrentTimestamp()
	return queueEndTime
}
