package userdb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"time"
)

var log = logrus.WithField("module", "entity/userdb")

var (
	// ErrClosed 已关闭.
	ErrClosed            = errors.New("ErrClosed")
	errInvalidShardIndex = errors.New("errInvalidShardIndex")
)

type useridToShard func(userid uint64) (int, error)

type FNoticeMsg func(userID uint64, noticeID int32, msg string, data []string)

type FMsgToService func(typ appframe.ServerType, msg proto.Message)

type FMsgRequest func(typ appframe.ServerType, msg proto.Message, cbk func(resp proto.Message, err error), timeout time.Duration)

// DB 分片管理
type DB struct {
	useridToShard useridToShard
	db            *sql.DB
	shards        []*shard
	wg            sync.WaitGroup
	closed        bool
	rw            sync.RWMutex

	msgToService FMsgToService
	msgRequest   FMsgRequest

	config *conf.GameConfig
	app    *appframe.Application

	AreaID int32

	Dev bool
}

// NewDB 创建 BD
func NewDB(db *sql.DB, shardCnt int, chanLen int, useridToShard useridToShard, msgToService FMsgToService, msgRequest FMsgRequest, app *appframe.Application, dev bool) *DB {
	d := new(DB)
	d.db = db
	d.useridToShard = useridToShard
	d.msgToService = msgToService
	d.msgRequest = msgRequest
	d.app = app
	d.Dev = dev

	d.shards = make([]*shard, shardCnt)
	d.wg.Add(shardCnt)
	for i := 0; i < shardCnt; i++ {
		s := newShard(d, chanLen, dev)
		d.shards[i] = s

		doWork := func() (done bool) {
			defer util.Recover()
			for f := range s.ch {
				f()
			}
			return true
		}

		util.SafeGo(func() {
			defer d.wg.Done()
			for !doWork() {
			}
		})
	}

	return d
}

func DBColumnValue(column *sql.ColumnType, pval interface{}) []string {
	var s_txt string
	switch v := pval.(type) {
	case nil:
		s_txt = "NULL"
	case time.Time:
		s_txt = v.Format("2006-01-02 15:04:05")
	case int, int8, int16, int32, int64, float32, float64, byte:
		s_txt = fmt.Sprint(v)
	case []byte:
		s_txt = fmt.Sprintf("%x", v)
	case bool:
		if v {
			s_txt = "1"
		} else {
			s_txt = "0"
		}
	default:
		s_txt = fmt.Sprint(v)
	}
	return []string{column.Name(), s_txt}
}

// GetUser 获取指定的用户, 回调函数 f 将会在独立的协程上执行.
func (d *DB) GetUser(userid uint64, f func(*User, error)) {
	idx, err := d.useridToShard(userid)
	if err != nil {
		f(nil, err)
		return
	}
	if idx < 0 || idx >= len(d.shards) {
		f(nil, errInvalidShardIndex)
		return
	}
	d.rw.RLock()
	if d.closed {
		d.rw.RUnlock()
		f(nil, ErrClosed)
		return
	}
	shard := d.shards[idx]
	shard.ch <- func() {
		u, err := shard.getUser(userid)
		if u != nil {
			u.config = shard.config
		}
		f(u, err)
	}
	d.rw.RUnlock()
}

// RecoverUser 获取指定的用户, 回调函数 f 将会在独立的协程上执行.
func (d *DB) RecoverUser(userid uint64, f func(*User, error)) {
	idx, err := d.useridToShard(userid)
	if err != nil {
		f(nil, err)
		return
	}
	if idx < 0 || idx >= len(d.shards) {
		f(nil, errInvalidShardIndex)
		return
	}
	d.rw.RLock()
	if d.closed {
		d.rw.RUnlock()
		f(nil, ErrClosed)
		return
	}
	shard := d.shards[idx]
	shard.ch <- func() {
		u, err := shard.recoverUser(userid)
		if u != nil {
			u.config = shard.config
		}
		f(u, err)
		shard.cache2DB(userid)
	}
	d.rw.RUnlock()
}

//不会阻塞当前线程，但可能会导致执行次序不准备问题
func (d *DB) SafeGetUser(userid uint64, f func(*User, error)) {
	idx, err := d.useridToShard(userid)
	if err != nil {
		f(nil, err)
		return
	}
	if idx < 0 || idx >= len(d.shards) {
		f(nil, errInvalidShardIndex)
		return
	}
	d.rw.RLock()
	if d.closed {
		d.rw.RUnlock()
		f(nil, ErrClosed)
		return
	}
	shard := d.shards[idx]
	postFun := func() {
		u, err := shard.getUser(userid)
		if u != nil {
			u.config = shard.config
		}
		f(u, err)
	}
	//如果插入不成功，就开协程插入
	select {
	case shard.ch <- postFun:
	default:
		util.SafeGo(func() {
			d.GetUser(userid, f)
		})
	}

	d.rw.RUnlock()
}

//share的负载，测试用
func (d *DB) ShareWorkerLen() int {
	var sum int
	for _, v := range d.shards {
		sum += len(v.ch)
	}
	return sum
}

func (d *DB) SharedCacheLen() int64 {
	var sum int64
	for _, v := range d.shards {
		sum += v.GetCacheInfo()
	}
	return sum
}

func (d *DB) EveryShardUpdateConf(config *conf.GameConfig) {
	d.config = config
	for _, v := range d.shards {
		share := v
		share.ch <- func() {
			share.config = config
		}
	}
}

// Close 关闭
func (d *DB) Close() error {
	d.rw.Lock()
	//defer d.rw.Unlock()
	d.closed = true
	for _, s := range d.shards {
		close(s.ch)
	}
	d.rw.Unlock()
	d.wg.Wait()
	var err error
	for _, s := range d.shards {
		e := s.sync()
		if e != nil {
			err = e
		}
	}
	return err
}

func (d *DB) loadUser(userid uint64) (*User, error) {
	u := NewUser(userid)
	//u.userid = userid
	u.db = d
	idx, err := d.useridToShard(userid)
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= len(d.shards) {
		return nil, err
	}
	u.config = d.shards[idx].config

	err = FillUser(d.db, u)
	if err != nil {
		return nil, err
	}

	u.checkInit()
	return u, nil
}

func FillUser(db *sql.DB, u *User) error {
	err := fillAccountInfo(db, u)
	if err != nil {
		return err
	}
	return fillUserInfo(db, u)
}

func fillAccountInfo(db *sql.DB, u *User) error {
	userId := u.userid
	err := db.QueryRow("SELECT account,account_type,loginid,nickname,head_img_url,login_ip,sex,status,created_ip,created_at FROM account WHERE userid=?", userId).Scan(
		&u.account, &u.accountType, &u.thirdAccountId, &u.nickname, &u.headImgUrl, &u.loginIP, &u.sex, &u.status, &u.createIp, &u.createTime)
	if err != nil {
		log.WithError(err).WithField("userid", userId).Error("LoadUserFailed Query")
		stack := debug.Stack()
		logrus.WithFields(logrus.Fields{
			"err":    err,
			"stack":  string(stack),
			"userId": userId,
		}).Error("loadUserFailed")

		os.Stderr.Write([]byte(fmt.Sprintf("%v\n", err)))
		os.Stderr.Write(stack)
		return err
	}
	return nil
}

func fillUserInfo(db *sql.DB, u *User) error {
	userId := u.userid

	var result []interface{}
	prefix := "SELECT level,exp,login_time,logout_time"
	result = append(result, &u.blobData.Level, &u.blobData.Exp, &u.blobData.LoginTime, &u.blobData.LogoutTime)

	var middle string
	composData := make(map[string]*[]byte, 0)
	for _, v := range u.components {
		filedName := v.dbFieldName()
		middle += fmt.Sprintf(",%s ", v.dbFieldName())
		data := make([]byte, 0)
		composData[filedName] = &data
		result = append(result, &data)
	}

	userTableName := "user"
	suffix := fmt.Sprintf(" FROM %s WHERE userid=?", userTableName)

	sql := prefix + middle + suffix

	err := db.QueryRow(sql, userId).Scan(result...)

	if err != nil {
		log.WithError(err).WithField("userid", userId).Error("LoadUserFailed Scan")
		return err
	}

	u.level = u.blobData.Level
	u.exp = u.blobData.Exp
	u.loginTime = u.blobData.LoginTime
	u.logoutTime = u.blobData.LogoutTime

	for _, component := range u.components {
		dbFieldName := component.dbFieldName()
		err := component.init(u, *composData[dbFieldName])
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{
				"userid":      u.userid,
				"dbFieldName": dbFieldName,
			}).Error("user component init")
			return err
		}
	}
	return nil
}

func (u *User) checkInit() {
	if u.level == 0 {
		logrus.WithField("userid", u.ID()).Warn("checkInit user level ==0")
		u.SetLevel(1)
	}
}

func (d *DB) updateUser(userid uint64, dirty map[string]interface{}) error {
	var fields = make([]string, 0, len(dirty))
	var values = make([]interface{}, 0, len(dirty))
	for k, v := range dirty {
		fields = append(fields, k+"=?")
		values = append(values, v)
	}
	userTableName := "user"
	cmd := fmt.Sprintf("UPDATE %s SET %s WHERE userid=?", userTableName, strings.Join(fields, ","))
	args := append(values, userid)
	_, err := d.db.Exec(cmd, args...)
	if err != nil {
		log.WithError(err).WithField("userid", userid).Error("Update user data failed")
		return err
	}
	return nil
}

// Sync 同步数据到 DB. TODO dirty处理方式要优化
func (u *User) sync() error {
	for _, v := range u.components {
		dbFieldName := v.dbFieldName()
		if v.isDirty() {
			data, err := proto.Marshal(v.toProtoMessage())
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"userID":      u.ID(),
					"dbFieldName": dbFieldName,
				}).WithError(err).Error("Sync error")
				continue
			}
			u.dirty[dbFieldName] = data
		}
	}

	if len(u.dirty) == 0 {
		return nil
	}

	err := u.db.updateUser(u.userid, u.dirty)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": u.ID(),
		}).WithError(err).Error("Sync updateUser error")
		return err
	}
	//数据库更新成功再清dirty
	u.dirty = make(map[string]interface{})
	for _, v := range u.components {
		v.clearDirty()
	}
	return nil
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
//表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}
