package userdb

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/proto/db"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"time"

	"github.com/golang/protobuf/proto"

	"sanguosha.com/sgs_herox/gameutil"
)

type char struct {
	BaseComp

	user     *User
	charInfo db.CharInfo
}

func (p *char) dbFieldName() string {
	return "char_info"
}

func (p *char) init(user *User, data []byte) error {
	p.user = user

	var dbChar db.CharInfo
	err := proto.Unmarshal(data, &dbChar)
	if err != nil {
		return err
	}

	p.charInfo = dbChar

	p.checkInit()
	return nil
}

func (p *char) toProtoMessage() proto.Message {
	return &p.charInfo
}

func (p *char) onLogin() {
	p.checkDayReset()
}

func (p *char) onLogout() {
	loginTime := p.user.loginTime.Unix()
	now := time.Now().Unix()

	loginZeroTime := gameutil.GetTargetDaySecond2ZeroTimeInt(loginTime)
	todayZeroTime := gameutil.GetTargetDaySecond2ZeroTimeInt(now)

	//登入登出不同天，取此时至0点时间为当日在线时间
	if loginZeroTime == todayZeroTime {
		p.charInfo.OnlineSecToday += int32(now - loginTime)
	} else {
		p.charInfo.OnlineSecToday = int32(now - todayZeroTime)
	}
	p.charInfo.OnlineSecTotal += (now - loginTime)
	p.setDirty()
}

func (p *char) onClock0() {
	p.checkDayReset()
}

func (p *char) onMinute() {

}

func (p *char) checkInit() {
	if p.charInfo.Init {
		return
	}

	p.charInfo.Init = true

	p.setDirty()
}

func (u *User) GetCharInfo() *db.CharInfo {
	return &u.Char.charInfo
}

func (p *char) checkDayReset() {
	if p.charInfo.LastResetTime >= gameutil.GetTodayZeroTimeInt() {
		return
	}
	lastResetTime := p.charInfo.LastResetTime
	p.charInfo.LastResetTime = time.Now().Unix()
	p.setDirty()

	lastResetZeroTime := gameutil.GetTargetDaySecond2ZeroTimeInt(lastResetTime)
	nowDayZeroTime := gameutil.GetTargetDaySecond2ZeroTimeInt(p.charInfo.LastResetTime)

	continuous := true
	if nowDayZeroTime-lastResetZeroTime > 24*3600 {
		continuous = false
	}

	p.dayReset(continuous)
}

func (p *char) dayReset(continuous bool) {
	if continuous {
		p.charInfo.C11LoginDayCnt++
	} else {
		p.charInfo.C11LoginDayCnt = 0
	}
	p.charInfo.LoginDayCnt++
	p.charInfo.OnlineSecToday = 0

	p.charInfo.CurDayWinCountText = 0
	p.charInfo.CurDayLoseCountText = 0
	p.charInfo.CurDayWinCountVoice = 0
	p.charInfo.CurDayLoseCountVoice = 0

	p.setDirty()
}

func (p *char) SetLoginDayCnt(v int32) {
	p.charInfo.LoginDayCnt = v
	p.setDirty()
}

func (p *char) ClearOnlineSecToday() {
	p.charInfo.OnlineSecToday = 0
	p.setDirty()
}
func (p *char) AddCardGroup(cg *gamedef.CardGroup) {
	cnt := len(p.charInfo.CardGroups)
	cg.GroupId = int32(cnt) + 1
	p.charInfo.CardGroups = append(p.charInfo.CardGroups, cg)
	p.charInfo.NowUseId = cg.GroupId
	p.setDirty()
}
func (p *char) ModifyCardGroup(id int, cg *gamedef.CardGroup) {
	p.charInfo.CardGroups[id] = cg
	p.charInfo.NowUseId = p.charInfo.CardGroups[id].GroupId
	p.setDirty()
}
func (p *char) SelectCardGroup(id int32) {
	p.charInfo.NowUseId = id
	p.setDirty()
}

func (p *char) AddUpWinLoseCount(gameMode gameconf.GameModeTyp, winLoseTyp gameconf.WinLoseTyp) {
	switch gameMode {
	case gameconf.GameModeTyp_MGTSpyText:
		switch winLoseTyp {
		case gameconf.WinLoseTyp_WLTWin:
			p.charInfo.WinCountText++
			p.charInfo.CurDayWinCountText++
		case gameconf.WinLoseTyp_WLTLose:
			p.charInfo.LoseCountText++
			p.charInfo.CurDayLoseCountText++
		default:

		}
	case gameconf.GameModeTyp_MGTSpyVoice:
		switch winLoseTyp {
		case gameconf.WinLoseTyp_WLTWin:
			p.charInfo.WinCountVoice++
			p.charInfo.CurDayWinCountVoice++
		case gameconf.WinLoseTyp_WLTLose:
			p.charInfo.LoseCountVoice++
			p.charInfo.CurDayLoseCountVoice++
		default:

		}
	default:

	}
	p.setDirty()
}

func (p *char) AddScore(score []*db.Score) {
	if p.charInfo.Score == nil {
		p.charInfo.Score = []*db.Score{}
		p.charInfo.Score = append(p.charInfo.Score, score...)
		return
	}

	for i := 0; i < len(p.charInfo.Score); i++ {
		for _, scoreList := range score {
			if scoreList.ScoreType == p.charInfo.Score[i].ScoreType {
				updateScore(&p.charInfo.Score[i].Score, scoreList.Score)
			}
		}
	}
	logrus.Debugf("user id %v,update hide score %v", p.user.userid, p.charInfo.Score)
	p.setDirty()
}

func updateScore(dbScore *[]*db.Int32KV, addScore []*db.Int32KV) {
	if len(addScore) < 1 {
		return
	}
	if !gameutil.IsInKVSliceK(addScore[0].Key, *dbScore) {
		*dbScore = append(*dbScore, addScore[0])
	} else {
		for j := 0; j < len(*dbScore); j++ {
			if (*dbScore)[j].Key == addScore[0].Key {
				if (*dbScore)[j].V != 0 {
					(*dbScore)[j].V += addScore[0].V
					(*dbScore)[j].V /= 2
				}
				break
			}
		}
	}
}
