package userdb

import (
	"errors"
	"fmt"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"

	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sanguosha.com/sgs_herox/gameutil"
)

var (
	// ErrIllegal 非法的值
	ErrIllegal = errors.New("ErrIllegal")
	// ErrDuplicate 重复的值
	ErrDuplicate = errors.New("ErrDuplicate")
	//货币不足
	ErrCurrencyInsufficient = errors.New("ErrCurrentInsufficient")
	//物品不存在
	ErrProductNotHave = errors.New("ErrProductNotHave")
)

type noCopy struct{}

func (*noCopy) Lock() {}

type UserStatus int

const (
	Offline UserStatus = iota
	Online
)

//这里的field和数据库里的字段相同,也和dirty键值相同，用于检测有数据改变但未标记为脏的情况 见CheckAndsetDirty
type UserDBData struct {
	Level int32
	Exp   int32

	LoginTime  time.Time //登入时间
	LogoutTime time.Time //登入时间//上次登出时间
}

// User 用户数据对象
type User struct {
	//data from db
	userid          uint64
	accountType     int32
	account         string
	nickname        string
	sex             int32
	headImgUrl      string
	headFrameImgUrl string
	thirdAccountId  string
	createIp        string

	level    int32
	exp      int32
	blobData UserDBData //玩家的数据库原始数据

	loginTime, logoutTime, createTime time.Time //登入时间
	bornTime                          int64     //出生时间，第一次初始化名字的时间

	//db
	dirty map[string]interface{}
	db    *DB

	//memory data
	loginIP        string
	status         UserStatus
	apkVersion     string
	apkFullVersion string
	config         *conf.GameConfig

	//send session
	sendReady bool //当玩家ReqMyData之前，不允许发消息给客户端

	sessionID appframe.SessionID
	Token     string
	TokenData string

	//components
	Char char
	Prop prop

	components []Component
	noCopy
}

func NewUser(userId uint64) *User {
	u := new(User)
	u.userid = userId
	u.dirty = make(map[string]interface{})
	u.status = Offline
	u.initComponent()
	return u
}
func (u *User) GetSessionID() appframe.SessionID {
	return u.sessionID
}
func (u *User) initComponent() {
	u.addComponent(&u.Char)
	u.addComponent(&u.Prop)
}

func (u *User) addComponent(component Component) {
	u.components = append(u.components, component)
}

func (u *User) Config() *conf.GameConfig {
	return u.config
}

// ID 用户标识符
func (u *User) ID() uint64 {
	return u.userid
}

func (u *User) AccountType() gameconf.AccountLoginTyp {
	return gameconf.AccountLoginTyp(u.accountType)
}

func (u *User) HeadImgUrl() string {
	return u.headImgUrl
}

func (u *User) BornTime() int64 {
	return u.bornTime
}

// Nickname 用户昵称
func (u *User) Nickname() string {
	return u.nickname
}

func (u *User) Account() string {
	return u.account
}

func (u *User) CreatedIP() string {
	return u.createIp
}

func (u *User) Sex() int32 {
	return u.sex
}

func (u *User) Level() int32 {
	return u.level
}

func (u *User) SetLevel(v int32) {
	u.level = v
	u.dirty["level"] = v
}

func (u *User) SetExp(v int32) {
	u.exp = v
	u.dirty["exp"] = v
}

func (u *User) Exp() int32 {
	return u.exp
}

func (u *User) GetOnlineSecToday() int32 {
	loginInTime := u.LoginTime()

	loginTime := loginInTime.Unix()
	now := time.Now().Unix()

	loginZeroTime := gameutil.GetTargetDaySecond2ZeroTimeInt(loginTime)
	todayZeroTime := gameutil.GetTargetDaySecond2ZeroTimeInt(now)

	if u.IsOnline() {
		if loginZeroTime != todayZeroTime {
			return int32(now - todayZeroTime)
		} else {
			return u.GetCharInfo().OnlineSecToday + int32(now-loginTime)
		}
	} else {
		if loginZeroTime != todayZeroTime {
			return 0
		} else {
			return u.GetCharInfo().OnlineSecToday
		}
	}
}

func (u *User) GetOnlineSecTotal() int64 {
	loginInTime := u.LoginTime()
	if u.IsOnline() {
		return u.GetCharInfo().OnlineSecTotal + gameutil.GetCurrentTimestamp() - loginInTime.Unix()
	} else {
		return u.GetCharInfo().OnlineSecTotal
	}
}

func (u *User) SetVersion(version string) {
	if version == "" {
		logrus.Error("version is empty")
	}
	ss := strings.Split(version, ";")
	if len(ss) != 2 {
		return
	}
	u.apkVersion = ss[0]
	u.apkFullVersion = version
}

func (u *User) GetApkFullVersion() string {
	return u.apkFullVersion
}

func (u *User) GetApkVersion() string {
	if u.apkVersion == "" {
		return "2.0.2"
	}
	return u.apkVersion
}

func (u *User) SetLoginIP(ip string) {
	u.loginIP = ip
}

func (u *User) LoginIP() string {
	return u.loginIP
}

func (u *User) SetNickname(nickname string) {
	u.nickname = nickname
}

func (u *User) SetHeadImgUrl(headImgUrl string) {
	u.headImgUrl = headImgUrl
}

func (u *User) SetHeadFrameImgUrl(headFrameImgUrl string) {
	u.headFrameImgUrl = headFrameImgUrl
}

func (u *User) SetThirdAccountId(v string) {
	u.thirdAccountId = v
}

func (u *User) ThirdAccountId() string {
	return u.thirdAccountId
}

func (u *User) SetSex(sex int32) {
	u.sex = sex
}

// LoginTime 登陆时间
func (u *User) LoginTime() time.Time {
	return u.loginTime
}

// LoginTime 登出时间
func (u *User) LogoutTime() time.Time {
	return u.logoutTime
}

func (u *User) CreatedTime() time.Time {
	return u.createTime
}

// 登入
func (u *User) SetLoginTime(t time.Time) {
	u.loginTime = t
	u.dirty["login_time"] = u.loginTime
}

func (u *User) SetOnline(online bool) {
	if online {
		u.status = Online
	} else {
		u.status = Offline
	}
}

func (u *User) Status() int32 {
	return int32(u.status)
}

func (u *User) SetSendReady(v bool) {
	u.sendReady = v
}

func (u *User) IsSendReady() bool {
	return u.sendReady
}

func (u *User) Login(ip string, auth *gamedef.AuthInfo) {
	u.SetNickname(auth.Nickname)
	u.SetHeadImgUrl(auth.HeadImgUrl)
	u.SetHeadFrameImgUrl(auth.HeadFrameImgUrl)
	u.SetSex(auth.Sex)
	u.SetThirdAccountId(auth.ThirdAccountId)
	u.SetLoginTime(time.Now())
	u.SetLoginIP(ip)
	u.SetOnline(true)
	u.OnLoginEvent()
}

func (u *User) Logout() {
	u.SetLogoutTime()
	u.SetOnline(false)
	u.SetSendReady(false)
	u.onLogoutEvent()
}

//登出
func (u *User) SetLogoutTime() {
	u.logoutTime = time.Now()
	u.dirty["logout_time"] = u.logoutTime
	logrus.Debugf("SetLogoutTime %d %s", u.userid, gameutil.ParseTimestamp2String(u.logoutTime.Unix()))
}

func (u *User) SendToClient(msg proto.Message) error {
	if !u.IsSendReady() {
		return nil
	}
	if u.sessionID.ID <= 0 {
		err := fmt.Errorf("send message to client error")
		logrus.WithFields(logrus.Fields{
			"sessionID": u.sessionID.ID,
			"serverID":  u.sessionID.SvrID,
		}).Debug(err)
		return err
	}
	u.db.app.GetSession(u.sessionID).SendMsg(msg)
	return nil
}

func (u *User) SendDebugInfo(msg string) {
	if !u.db.Dev {
		return
	}
	u.SendToClient(&cmsg.SNoticeDebugInfo{Msg: msg})
}

//绑定session
func (u *User) BindSession(id appframe.SessionID) {
	u.sessionID = id
}

//解除绑定
func (u *User) UnbindSession() {
	u.sessionID = appframe.SessionID{}
}

func (u *User) ServerNotice(noticeID int32, msg string, data []string) {
	app := u.db.app
	notice := &smsg.ServerNotice{AppID: app.ID(), UserID: u.ID(), NoticeID: noticeID, Msg: msg, Data: data}

	serverList := app.GetAvailableServerIDs(sgs_herox.SvrTypeGate)
	for _, id := range serverList {
		app.GetServer(id).SendMsg(notice)
	}
}

//是否在线
func (u *User) IsOnline() bool {
	return u.status == Online
}

func (u *User) OnMinuteEvent() {
	if !u.IsOnline() {
		return
	}

	for _, v := range u.components {
		v.onMinute()
	}
}

func (u *User) OnZero0ClockEvent(now time.Time) {
	if !u.IsOnline() {
		return
	}
	for _, v := range u.components {
		v.onClock0()
	}
}

func (u *User) OnLoginEvent() {
	for _, v := range u.components {
		v.onLogin()
	}
}

func (u *User) onLogoutEvent() {
	for _, v := range u.components {
		v.onLogout()
	}
}

func (u *User) SendToLobby(msg proto.Message) {
	u.db.msgToService(sgs_herox.SvrTypeLobby, msg)
}

//尝试向本协程插入事件，如果插入不进去，开启goroutine进行插入
func (u *User) SafePost(fun func()) {
	u.db.SafeGetUser(u.ID(), func(user *User, err error) {
		if err != nil {
			return
		}
		fun()
	})
}

func (u *User) CharLevel() int32 {
	return u.level
}
