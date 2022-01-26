package entity

import (
	"time"

	"github.com/sirupsen/logrus"

	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/entity/userdb"
	"sanguosha.com/sgs_herox/gameutil"
)

const AM0FreshIntervel = 20 * time.Millisecond //分批刷新的时间间隔
const AM0FreshCntPer = 500                     //一次刷新的数量

//----------------------------------定时任务全放这里-------------------------------------
//0点刷新
func AM0Fresh() {
	logrus.Info("AM0Fresh Come")
	onlineUserList := make([]uint64, 0)
	for _, userID := range onlineMgr.AllUserIds() {
		onlineUserList = append(onlineUserList, userID)
	}
	AM0FreshBatch(onlineUserList)
}

func AM0FreshBatch(userIDList []uint64) {
	util.SafeGo(func() {
		//开启个线程，分批执行0点刷新功能
		ticker := time.NewTicker(AM0FreshIntervel)
		defer func() {
			ticker.Stop()
			logrus.Info("AM0FreshBatch Over")
		}()
		now := time.Now()
		logrus.WithFields(logrus.Fields{
			"time": now,
		}).Info("AM0FreshBatch Start")

		for {
			select {
			case <-ticker.C:
				freshCnt := 0
				for _, userID := range userIDList {
					UserDBInstance.GetUser(userID, func(u *userdb.User, err error) {
						if err != nil {
							return
						}
						u.OnZero0ClockEvent(now)
					})
					freshCnt++
					if freshCnt >= AM0FreshCntPer {
						break
					}
				}
				userIDList = userIDList[freshCnt:]
				if len(userIDList) == 0 {
					return
				}
			}
		}
	})
}

//每分钟刷新
func EveryMinuteFresh() {
	onlineUserList := make([]uint64, 0)
	for _, userID := range onlineMgr.AllUserIds() {
		onlineUserList = append(onlineUserList, userID)
	}
	EveryMinuteFreshBatch(onlineUserList)
}

func EveryMinuteFreshBatch(userIDList []uint64) {
	length := len(userIDList)
	if length == 0 {
		return
	}

	freshPerOnce := 100 //每次100个
	//计算下interval,由于是1分钟刷新，任务要在30秒内完成
	interval := time.Duration(gameutil.SafeDivFloat64(float64(30*time.Second), float64(length))) * time.Duration(freshPerOnce)

	if interval > time.Second || interval == 0 {
		interval = time.Second
	}
	util.SafeGo(func() {
		//开启个线程，分批执行刷新功能
		ticker := time.NewTicker(interval)
		defer func() {
			ticker.Stop()
			//logrus.Info("EveryMinuteFreshBatch Over")
		}()
		for {
			select {
			case <-ticker.C:
				freshCnt := 0
				for _, userID := range userIDList {
					UserDBInstance.GetUser(userID, func(u *userdb.User, err error) {
						if err != nil {
							return
						}
						u.OnMinuteEvent()
					})
					freshCnt++
					if freshCnt >= freshPerOnce {
						break
					}
				}
				userIDList = userIDList[freshCnt:]
				if len(userIDList) == 0 {
					return
				}
			}
		}
	})
}
