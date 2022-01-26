package lobby

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameutil"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"time"
)

const ROOM_ID_BEGIN = 1000

type RoomManager struct {
	IdMaker uint32

	roomId2room map[uint32]*Room
	roomNO2room map[uint32]*Room
}

func NewRoomManager() *RoomManager {
	p := &RoomManager{}
	p.init()
	return p
}

func (p *RoomManager) init() {
	p.roomId2room = make(map[uint32]*Room)
	p.roomNO2room = make(map[uint32]*Room)
}

func (p *RoomManager) newRoomId() (uint32, bool) {
	p.IdMaker++
	if p.IdMaker < ROOM_ID_BEGIN {
		p.IdMaker = ROOM_ID_BEGIN
	}

	if _, exist := p.roomId2room[p.IdMaker]; !exist {
		return p.IdMaker, true
	}

	return 0, false
}

// 生成队伍编号
func (p *RoomManager) newRoomNO() (uint32, bool) {
	for i := 0; i < 100; i++ {
		roomNo := uint32(gameutil.RandomBetween(100000, 999999))

		if !p.isRoomNOExist(roomNo) {
			return roomNo, true
		}
	}
	return 0, false
}

func (p *RoomManager) isRoomNOExist(teamNo uint32) bool {
	_, exist := p.roomNO2room[teamNo]
	if !exist {
		return false
	}
	return true
}

func (p *RoomManager) defaultRoomSetting(gameMode gameconf.GameModeTyp) *gamedef.RoomSetting {
	return &gamedef.RoomSetting{
		GameMode:   gameMode,
		RoomName:   gameConfig.Global.RoomSettingDefaultRoomName,
		MaxPlayer:  gameConfig.Global.RoomSettingDefaultPlayerCount,
		AllowEnter: gameConfig.Global.RoomSettingDefaultAllowEnter,
	}
}

func (p *RoomManager) checkRoomCfg(setting *gamedef.RoomSetting) error {
	if setting == nil {
		return errors.New("setting ==nil")
	}
	//if setting.RoomName == "" || len(setting.RoomName) > 100 {
	//	return fmt.Errorf("roomName len:%v not allow", len(setting.RoomName))
	//}
	switch setting.GameMode {
	case gameconf.GameModeTyp_MGTSpyText:
	case gameconf.GameModeTyp_MGTSpyVoice:
	default:
		return fmt.Errorf("gameMode:%v not support", setting.GameMode)
	}

	return nil
}

func (p *RoomManager) newRoom(setting *gamedef.RoomSetting, creator uint64) (*Room, error) {
	err := p.checkRoomCfg(setting)
	if err != nil {
		return nil, err
	}

	roomId, ok := p.newRoomId()
	if !ok {
		return nil, errors.New("newRoomId gen error")
	}

	roomNO, ok := p.newRoomNO()
	if !ok {
		return nil, errors.New("newTeamNO error")
	}

	t := &Room{
		roomId:  roomId,
		roomNO:  roomNO,
		setting: setting,
		Owner:   creator,
		Seats:   make([]*RoomSeat, setting.MaxPlayer),
		lookers: make(map[uint64]*user, 0),
	}

	for i := uint32(0); i < setting.MaxPlayer; i++ {
		t.Seats[i] = new(RoomSeat)
	}

	t.voiceId = makeVoiceId(roomId)

	p.roomId2room[roomId] = t
	p.roomNO2room[roomNO] = t
	return t, nil
}

func makeVoiceId(roomId uint32) string {
	return fmt.Sprintf("%d_%d_%d", sgs_herox.AppID, roomId, time.Now().Unix())
}

func (p *RoomManager) delRoomById(roomId uint32) bool {
	team, ok := p.roomId2room[roomId]
	if !ok {
		logrus.WithFields(logrus.Fields{
			"roomId": roomId,
		}).Debug("delTableByID room not found")
		return false
	}

	delete(p.roomId2room, roomId)
	delete(p.roomNO2room, team.roomNO)
	team.clearUser()
	return true
}

func (p *RoomManager) findRoomById(roomId uint32) (*Room, bool) {
	team, ok := p.roomId2room[roomId]
	return team, ok
}

func (p *RoomManager) findRoomByRoomNO(roomNO uint32) (*Room, bool) {
	team, ok := p.roomNO2room[roomNO]
	return team, ok
}

func (p *RoomManager) findRoomCanEnter(gameMode gameconf.GameModeTyp) (*Room, bool) {
	for _, room := range p.roomId2room {
		if room.setting.GameMode != gameMode {
			continue
		}

		if room.canSeatNewPlayer() {
			return room, true
		}
	}
	return nil, false
}
