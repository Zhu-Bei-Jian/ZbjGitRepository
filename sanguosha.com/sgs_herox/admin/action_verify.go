package admin

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/gameshared/manager"
	"sanguosha.com/sgs_herox/gameutil"
	"sync"
	"time"
)

var avs map[uint64]*ActionVerify
var avsIdMgr *manager.IDManager
var avsLock sync.Mutex

type ActionVerify struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Who  string `json:"who"`
	Time string `json:"time"`

	Action func() `json:"-"`
}

func InitActionVerify(app *appframe.Application) {
	avsIdMgr = &manager.IDManager{}
	avsIdMgr.Init(app.ID())

	avs = make(map[uint64]*ActionVerify)
}

func AddActionVerify(who string, name string, desc string, f func()) {
	avsLock.Lock()
	defer avsLock.Unlock()
	id := avsIdMgr.GeneratePKID()

	info := &ActionVerify{Who: who, ID: id, Name: name, Desc: desc, Action: f, Time: time.Now().Format("2006-01-02 15:04:05")}
	avs[id] = info

	logrus.Info("新增审核:", id, " ", who, " ", name, " ", desc)
	gameutil.SafeCallAfter(300*time.Second, func() {
		avsLock.Lock()
		defer avsLock.Unlock()

		if _, ok := avs[id]; ok {
			delete(avs, id)
			logrus.Info("审核超时:", id, " ", who, " ", name, " ", desc)
		}
	})
}

func DoActionVerify(id uint64, t int32) bool {
	avsLock.Lock()
	defer avsLock.Unlock()

	if v, ok := avs[id]; ok {
		delete(avs, id)
		if t == 1 {
			//app.Post(func() {
			util.SafeGo(func() {
				v.Action()
			})
			//})
		}
		return true
	}
	return false
}

func GetVerifyInfo() string {
	avsLock.Lock()
	defer avsLock.Unlock()
	l := []*ActionVerify{}
	for _, v := range avs {
		l = append(l, v)
	}
	b, e := json.Marshal(l)
	if e != nil {
		//logrus.Debug("GetVerifyInfo:", e)
	}
	//logrus.Debug("GetVerifyInfo:", string(b))
	return string(b)
}
