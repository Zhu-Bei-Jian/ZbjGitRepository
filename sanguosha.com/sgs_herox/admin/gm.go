package admin

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/gameshared/accountservice"
	"sanguosha.com/sgs_herox/gameutil"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	reqCallTimeout = time.Second * 30
)

var ErrParamsWrong = errors.New("ErrParamsWrong")

type gmManager struct {
	account2userID map[string]uint64

	number    int64 // 操作批号
	startTime int64 // 启动时间

	ddClose chan bool
	ddList  []string
	ddLock  sync.Mutex
}

func NewGMManager() *gmManager {
	gm := &gmManager{}
	gm.Init()
	return gm
}

func (p *gmManager) Init() {
	p.account2userID = make(map[string]uint64)
	p.startTime = gameutil.GetCurrentTimestamp()
	p.ddClose = make(chan bool)
}

// 获取操作ID, 每次操作是唯一的
func (p *gmManager) getGMOperateID() string {
	atomic.AddInt64(&p.number, 1)
	return fmt.Sprintf("%d_%d", p.startTime, p.number)
}

func (p *gmManager) onGMCommand(params map[string]string, receiver commandCallback) {
	opID := p.getGMOperateID()
	command := params["command"]
	res, err := p.process(opID, command, params)
	// 记录文件日志
	logrus.WithFields(logrus.Fields{
		"cmd":         command,
		"res":         res,
		"err":         err,
		"GMOperateID": opID,
	}).Info("onGMCommand")
	receiver.writeMsg(res, err)
}

//注意：这个函数是同步执行
func (p *gmManager) onRemoteCommand(cmd string, callback func(s string)) {
	cmdList := strings.Split(cmd, " ")
	switch strings.ToLower(cmdList[0]) {
	case "/status":
		gmGetServerStatus(callback)
	case "/q":
		gmGetUserInfo(cmd, callback)
	case "/closeserver":
		if !appConfig.Develop {
			return
		}
		gmCloseServer(cmd, callback)
	case "/sendhorselamb":
		gmSendHorseLamp(cmd, callback)
	default:
		callback(fmt.Sprintf("收到GM命令无效: " + cmd))
	}
}

//部分暂未支持
func (p *gmManager) process(opID string, command string, params map[string]string) (interface{}, error) {
	subCommand := params["subCommand"]
	switch command {
	case "getAccountInfo":
		return p.onGetAccountInfo(subCommand, params["param0"])
	//case "getUserInfo":
	//	return queryUserDetailSync(subCommand, params["account"])
	case "horse_race_lamp":
		return p.onSendHorseRaceLamp(params)
	default:
		return nil, errors.New("命令不存在")
	}
	return nil, errors.New("命令不存在")
}

func (p *gmManager) onGetAccountInfo(subCommand string, param string) (interface{}, error) {
	var userIds []uint64
	switch subCommand {
	case "byAccount":
		userID, err := accountService.GetUserIDByAccountSync(param)
		if err != nil {
			if err == accountservice.ErrNotFound {
				return nil, errors.New("账号不存在")
			}
			return nil, err
		}
		userIds = append(userIds, userID)
	case "byNickname":
		userIdsTemp, err := accountService.GetUserIdsByNickNameSync(param, false)
		if err != nil {
			if err == accountservice.ErrNotFound {
				return nil, errors.New("昵称不存在")
			}
			return nil, err
		}
		if len(userIdsTemp) == 0 {
			return nil, fmt.Errorf("昵称:%v 不存在", param)
		}
		userIds = userIdsTemp
	case "byUserID":
		userID, err := gameutil.ToUInt64(param)
		if err != nil {
			return nil, fmt.Errorf("userID格式不对")
		}
		userIds = append(userIds, userID)
	default:
		return nil, errors.New("无效参数")
	}

	return nil, nil
}

func (p *gmManager) onSendHorseRaceLamp(paramMap map[string]string) (interface{}, error) {
	//message := paramMap["message"]
	//channelID, err := gameutil.ToInt32(paramMap["channelID"])
	//if err != nil {
	//	return nil, ErrParamsWrong
	//}
	//notice := &smsg.ServerNotice{AppID: app.ID(), NoticeID: 1, Msg: message, ChannelID: channelID}
	//app.GetService(sgs_herox.SvrTypeLobby).SendMsg(notice)
	return nil, nil
}
