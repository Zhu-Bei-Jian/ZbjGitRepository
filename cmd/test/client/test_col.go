package client

import (
	"reflect"
	"time"
)

type ColInfo struct {
	K string
	v []interface{}
}

type DelayInfo struct {
	total int64
	times int64
}

var infos chan *ColInfo

var login_success = 0
var login_err = 0
var last_login_success, last_login_err int

var game_start = 0
var game_over = 0
var last_game_start, last_game_over int

var pingDelays map[int32]*DelayInfo
var respDelays map[reflect.Type]*DelayInfo

func init() {
	infos = make(chan *ColInfo, 1000000)
	pingDelays = make(map[int32]*DelayInfo)
	respDelays = make(map[reflect.Type]*DelayInfo)
}

var last_log int64

func GoShowCol() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				pushCol("log", 1)
			}
		}
	}()
	go func() {
		for i := range infos {
			switch i.K {
			case "ping":
				if len(i.v) != 2 {
					return
				}
				svr := i.v[0].(int32)
				diff := i.v[1].(int64)
				if diff >= 0 {
					di, ok := pingDelays[svr]
					if !ok {
						di = &DelayInfo{
							total: diff,
							times: 1,
						}
						pingDelays[svr] = di
					} else {
						di.total += diff
						di.times += 1
					}
				}
			case "resp":
				if len(i.v) != 2 {
					return
				}
				refType := i.v[0].(reflect.Type)
				diff := i.v[1].(int64)
				if diff >= 0 {
					di, ok := respDelays[refType]
					if !ok {
						di = &DelayInfo{
							total: diff,
							times: 1,
						}
						respDelays[refType] = di
					} else {
						di.total += diff
						di.times += 1
					}
				}
			case "login_success":
				login_success++
			case "login_err":
				login_err++
			case "game_start":
				game_start++
			case "game_over":
				game_over++
			case "log":
				timeNow := time.Now().Unix()
				if (timeNow - last_log) >= 10 {
					last_log = timeNow
					for _, di := range pingDelays {
						if di.times != 0 {
							//Log("ping %v %d avg %d", appframe.ServerType(svr), di.times, di.total/1e6/di.times)
							di.total = 0
							di.times = 0
						}
					}
					for _, di := range respDelays {
						if di.times != 0 {
							//Log("%v %d avg %d", type2name(refType), di.times, di.total/1e6/di.times)
							di.total = 0
							di.times = 0
						}
					}
					if last_login_success != login_success || last_login_err != login_err {
						last_login_success = login_success
						last_login_err = login_err
						//logrus.Info(fmt.Sprintf("登录 成功数量 %d 失败数量 %d", login_success, login_err))
					}
					if last_game_start != game_start || last_game_over != game_over {
						last_game_start = game_start
						last_game_over = game_over
						//logrus.Info(fmt.Sprintf("游戏 开始数量 %d 结束数量 %d", game_start, game_over))
					}
				}
			default:
			}
		}
	}()
}
func pushCol(k string, v ...interface{}) {
	infos <- &ColInfo{k, v}
}
