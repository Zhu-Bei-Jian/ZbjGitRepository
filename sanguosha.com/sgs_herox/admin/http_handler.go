package admin

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"net/http"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/entity/userdb"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (p *httpServer) write(res map[string]interface{}, w http.ResponseWriter) {
	result, _ := json.MarshalIndent(res, "", "\t")
	logrus.WithFields(logrus.Fields{}).Debug("return ", string(result))
	w.Write([]byte(result))
}

func (p *httpServer) writeObj(res interface{}, w http.ResponseWriter) {
	result, _ := json.MarshalIndent(res, "", "\t")
	logrus.WithFields(logrus.Fields{}).Debug("return ", string(result))
	w.Write([]byte(result))
}

type WsArgMonitor struct {
	Addr        string
	ConnectInfo string
}

func (p *httpServer) gmCommand(w http.ResponseWriter, r *http.Request) {
	rcv := &receiver{
		res: make(map[string]interface{}),
		w:   w,
		r:   r,
	}

	r.ParseForm()
	params := make(map[string]string)
	for k, v := range r.Form {
		params[k] = v[0]
	}
	gmMgr.onGMCommand(params, rcv)
}

func (p *httpServer) reload(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	defer p.write(res, w)

	optType := r.FormValue("server")
	needReloadServerType := make([]appframe.ServerType, 0)
	switch optType {
	case "all":
		needReloadServerType = reloadServers
	case "auth":
		needReloadServerType = append(needReloadServerType, sgs_herox.SvrTypeAuth)
	case "lobby":
		needReloadServerType = append(needReloadServerType, sgs_herox.SvrTypeLobby)
	case "game":
		needReloadServerType = append(needReloadServerType, sgs_herox.SvrTypeGame)
	case "entity":
		needReloadServerType = append(needReloadServerType, sgs_herox.SvrTypeEntity)
	default:
		res["错误"] = "参数错误"
		return
	}
	//request all server，记录request servers,并等待全部返回
	var wg sync.WaitGroup
	allReloadServerIDs := make([]uint32, 0)
	reloadSuccessServerIDs := make([]uint32, 0)
	realoadFailServerIDs := make([]uint32, 0)
	for _, v := range needReloadServerType {
		serverIDs := app.GetAvailableServerIDs(v)
		for _, v := range serverIDs {
			wg.Add(1)
			serverID := v
			allReloadServerIDs = append(allReloadServerIDs, serverID)
			app.GetServer(serverID).ReqSugar(&smsg.AdAllReqReload{}, func(resp *smsg.AdAllRespReload, err error) {
				defer wg.Done()
				//TODO 记录成功和失败的
				if err != nil {
					realoadFailServerIDs = append(realoadFailServerIDs, serverID)
					return
				}

				if resp.ErrCode != smsg.AdAllRespReload_Invalid {
					realoadFailServerIDs = append(realoadFailServerIDs, serverID)
					return
				}
				reloadSuccessServerIDs = append(reloadSuccessServerIDs, serverID)
			}, time.Second*60)
		}
	}
	fmt.Println("allReloadServerID:", allReloadServerIDs)
	tmp := make(map[string]interface{})
	tmp["需要重载ServerID列表"] = allReloadServerIDs
	p.write(tmp, w)
	fmt.Fprint(w, "\n")
	wg.Wait()
	fmt.Printf("reloadSuccessServerIDs:%v,reloadFailServerIDs:%v", reloadSuccessServerIDs, realoadFailServerIDs)
	res["成功ServerID列表"] = reloadSuccessServerIDs
	res["失败ServerID列表"] = realoadFailServerIDs
	return
}

func closeServer(optType string, serverids []uint32) {
	closeServerIDs := make([]uint32, 0)
	switch optType {
	case "all":
		closeServerTypes := []appframe.ServerType{
			sgs_herox.SvrTypeGate,
			sgs_herox.SvrTypeGame,
			sgs_herox.SvrTypeAuth,
			sgs_herox.SvrTypeEntity,
			sgs_herox.SvrTypeLobby,
			sgs_herox.SvrTypeAccount,
		}
		for _, v := range closeServerTypes {
			if app.GetService(v).Available() {
				ids := app.GetAvailableServerIDs(v)
				closeServerIDs = append(closeServerIDs, ids...)
			}
		}
	//case "ai":
	//	//serverIDStr := r.FormValue("serverids")
	//	//serverIDs, err := gameutil.ParseUInt32s(serverIDStr)
	//	//if err != nil {
	//	//	res["error"] = "参数错误"
	//	//	return
	//	//}
	//
	//	ids := app.GetAvailableServerIDs(sgs_herox.SvrTypeAI)
	//	for _, id := range ids {
	//		for _, v := range serverids {
	//			if id == v {
	//				app.GetServer(id).SendMsg(&gameproto.ReqCloseServer{})
	//				//fmt.Fprintf(w, "已通知AI[%d] 关闭", v)
	//				fmt.Printf("已通知AI[%d] 关闭\n", v)
	//				break
	//			}
	//		}
	//	}
	//	return
	case "id":
		closeServerIDs = append(closeServerIDs, serverids...)
	case "gate":
		closeServerIDs = app.GetAvailableServerIDs(sgs_herox.SvrTypeGate)
	default:
		return
	}

	//var wg sync.WaitGroup
	successServerIDs := make([]uint32, 0)
	failServerIDs := make([]uint32, 0)

	fmt.Println("allNeedCloseServerIDs:", closeServerIDs)
	tmp := make(map[string]interface{})
	tmp["需要关闭的ServerID列表"] = closeServerIDs

	for _, v := range closeServerIDs {
		serverID := v
		var resp smsg.AdAllRespCloseServer
		err := app.GetServer(serverID).CallSugar(&smsg.AdAllReqCloseServer{}, &resp, time.Minute)

		if err != nil {
			failServerIDs = append(failServerIDs, serverID)
			fmt.Println("关闭失败", serverID, err)
			continue
		}

		if resp.ErrCode != smsg.AdAllRespCloseServer_Invalid {
			failServerIDs = append(failServerIDs, serverID)
			fmt.Println("关闭失败", serverID, resp.ErrCode)
			continue
		}
		successServerIDs = append(successServerIDs, serverID)
		fmt.Println("关闭成功", serverID)
	}

	fmt.Printf("SuccessServerIDs:%v,FailServerIDs:%v", successServerIDs, failServerIDs)
	if optType == "all" {
		gameutil.SafeCallAfter(15*time.Second, app.Exit)
		//app.Exit()
	}

}

func (p *httpServer) closeServer(w http.ResponseWriter, r *http.Request) {
	optType := r.FormValue("optType")
	serverIDStr := r.FormValue("serverids")
	serverIDs, _ := gameutil.ParseUInt32s(serverIDStr, ",")

	var desc string
	switch optType {
	case "all":
		desc = "关所有服（除了AI)"
	case "gate":
		desc = "关gate服"
	case "ai":
		desc = "关AI服"
	default:
		bs, _ := MakeCommonResult(0, "参数错误", nil)
		w.Write(bs)
	}

	AddActionVerify("a", "CloseServer", desc, func() {
		closeServer(optType, serverIDs)
	})
	bs, _ := MakeCommonResult(0, "已提交审核", nil)
	w.Write(bs)
	return
}

//func (p *httpServer) getRankInfo(w http.ResponseWriter, r *http.Request) {
//	res := make(map[string]interface{})
//	defer p.write(res, w)
//
//	rankType := r.FormValue("rankType")
//
//	switch rankType {
//	case "qualify":
//		seasonIDStr := r.FormValue("seasonID")
//		seasonID, err := gameutil.ToInt32(seasonIDStr)
//		if err != nil {
//			res["error"] = "参数错误"
//			return
//		}
//		rankList, err := shareddata.GetQualifyRankListNoLimit(seasonID)
//		type UserQualifyInfo struct {
//			UserID uint64
//			Star   int32
//		}
//		retList := make([]*UserQualifyInfo, 0)
//		for _, v := range rankList {
//			star := int32(v.Score / shareddata.RANK_SCORE_VALUE_TAG)
//			if star <= 0 {
//				continue
//			}
//			retList = append(retList, &UserQualifyInfo{
//				UserID: v.UserID,
//				Star:   star,
//			})
//		}
//		res["data"] = retList
//	default:
//		res["error"] = "参数错误"
//		return
//	}
//
//	return
//}

func (p *httpServer) ServerConfig(w http.ResponseWriter, r *http.Request) {
	t, e := ParseHtmlFiles("web/template/config.html")
	if e != nil {
		fmt.Print("ParseFiles:", e)
		return
	}
	data := make(map[string]string)

	infoMap, e := dbMgr.GetServerInfo()
	if e != nil {
	}
	if infoMap == nil {
		infoMap = make(map[string]string)
	}
	serverInfo, _ := json.MarshalIndent(infoMap, "", "\t")

	//openInfo := GetServerOpenInfo()
	if len(infoMap) == 0 {
		data["ServerOpenTime"] = "未知"
		data["ServerCloseTime"] = "未知"
		data["InnerIp"] = "未知"
		data["ServerInfo"] = string(serverInfo)
	} else {
		data["ServerOpenTime"], _ = infoMap["ServerOpenTime"]
		data["ServerCloseTime"], _ = infoMap["ServerCloseTime"]
		data["InnerIp"], _ = infoMap["InnerIp"]
		data["ServerInfo"] = string(serverInfo)
	}
	e = t.Execute(w, data)
	if e != nil {
		fmt.Print("Execute:", e)
		return
	}
}
func (p *httpServer) EditServerOpenTime(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("EditServerOpenTime")
	infoMap, err := dbMgr.GetServerInfo()
	if err != nil {
		bs, _ := MakeCommonResult(1, fmt.Sprintf("错误：%s", err.Error()), nil)
		w.Write(bs)
		return
	}

	r.ParseForm()
	info := r.FormValue("info")
	if len(info) != 0 {
		_, e := gameutil.ParseDatetime2Timestamp(info)
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
	}

	_, _, err = CheckServerInfo(infoMap, "ServerOpenTime", info)
	if err != nil {
		logrus.Error("EditServerOpenTime:", err)
		bs, _ := MakeCommonResult(0, err.Error(), nil)
		w.Write(bs)
		return
	}
	logrus.Info("EditServerOpenTime:", info)
	bs, _ := MakeCommonResult(0, "已提交", nil)
	w.Write(bs)
}
func (p *httpServer) EditServerCloseTime(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("EditServerCloseTime")

	infoMap, err := dbMgr.GetServerInfo()
	if err != nil {
		bs, _ := MakeCommonResult(1, fmt.Sprintf("错误：%s", err.Error()), nil)
		w.Write(bs)
		return
	}

	r.ParseForm()
	info := r.FormValue("info")

	t := int64(0)
	if len(info) != 0 {
		var e error
		t, e = gameutil.ParseDatetime2Timestamp(info)
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
	}

	if t != 0 {
		old := int64(0)
		if oldValue, ok := infoMap["ServerOpenTime"]; ok {
			old, _ = gameutil.ParseDatetime2Timestamp(oldValue)
		}
		if t < old {
			bs, _ := MakeCommonResult(1, "小于开服时间不允许设置，可选设置为空", nil)
			w.Write(bs)
			return
		}
	}

	_, _, err = CheckServerInfo(infoMap, "ServerCloseTime", info)
	if err != nil {
		logrus.Error("EditServerCloseTime:", err)
		bs, _ := MakeCommonResult(0, err.Error(), nil)
		w.Write(bs)
		return
	}
	logrus.Info("EditServerCloseTime:", info)
	bs, _ := MakeCommonResult(0, "已提交", nil)
	w.Write(bs)
}
func (p *httpServer) EditInnerIp(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("EditInnerIp")

	infoMap, err := dbMgr.GetServerInfo()
	if err != nil {
		bs, _ := MakeCommonResult(1, fmt.Sprintf("错误：%s", err.Error()), nil)
		w.Write(bs)
		return
	}

	r.ParseForm()
	info := r.FormValue("info")

	_, _, err = CheckServerInfo(infoMap, "InnerIp", info)
	if err != nil {
		logrus.Error("EditInnerIp:", err)
		bs, _ := MakeCommonResult(0, err.Error(), nil)
		w.Write(bs)
		return
	}
	logrus.Info("EditInnerIp:", info)
	bs, _ := MakeCommonResult(0, "已提交", nil)
	w.Write(bs)
}

func (p *httpServer) ServerCtrl(w http.ResponseWriter, r *http.Request) {
	t, e := ParseHtmlFiles("web/template/ctrl.html")
	if e != nil {
		fmt.Print("ParseFiles:", e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Print("Execute:", e)
		return
	}
}

func (p *httpServer) ActionVerify(w http.ResponseWriter, r *http.Request) {
	t, e := ParseHtmlFiles("web/template/action_verify.html")
	if e != nil {
		fmt.Print("ParseFiles:", e)
		return
	}
	data := make(map[string]interface{})

	data["Info"] = GetVerifyInfo()
	e = t.Execute(w, data)
	if e != nil {
		fmt.Print("Execute:", e)
		return
	}
}
func (p *httpServer) DoActionVerify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vid := r.FormValue("id")
	vt := r.FormValue("t")
	logrus.Debug("DoActionVerify:", vid, " ", vt)

	id, _ := strconv.ParseInt(vid, 10, 64)
	t, _ := strconv.Atoi(vt)

	if !DoActionVerify(uint64(id), int32(t)) {
		p.writeObj(&CommonResult{1, "审核项已不存在", nil}, w)
	} else {
		p.writeObj(&CommonResult{0, "", nil}, w)
	}
}
func SendToAllEntity(msg proto.Message) {
	ids := app.GetAvailableServerIDs(sgs_herox.SvrTypeEntity)
	for _, id := range ids {
		s := app.GetServer(id)
		s.SendMsg(msg)
	}
}
func (p *httpServer) KickAll(w http.ResponseWriter, r *http.Request) {
	//logrus.Debug("KickAll")

	AddActionVerify("a", "KickAll", "所有人下线", func() {
		msg := &smsg.AllGaNtfKickUserOut{
			KickAll: true,
			Reason:  gameconf.KickUserOutReason_KUOServerUpgrade,
		}
		gateids := app.GetAvailableServerIDs(sgs_herox.SvrTypeGate)
		for _, gateid := range gateids {
			app.GetServer(gateid).SendMsg(msg)
		}
	})
	//err := app.GetService(sgs_herox.SvrTypeLobby).SendMsg(&smsg.AllLoReqKickUserOut{
	//	UserID:    0,
	//	AllAction: 1,
	//})
	//if err != nil {
	//	logrus.Error("KickAll:", err)
	//	bs, _ := MakeCommonResult(1, err.Error(), nil)
	//	w.Write(bs)
	//	return
	//}
	bs, _ := MakeCommonResult(0, "已提交审核", nil)
	w.Write(bs)
}

func (p *httpServer) KickAllAndForceGameOver(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Silence")

	AddActionVerify("a", "KickAllAndForceGameOver", "所有人下线，结束所有游戏", func() {
		msg := &smsg.AllGaNtfKickUserOut{
			KickAll: true,
			Reason:  gameconf.KickUserOutReason_KUOServerUpgrade,
		}
		gateids := app.GetAvailableServerIDs(sgs_herox.SvrTypeGate)
		for _, gateid := range gateids {
			app.GetServer(gateid).SendMsg(msg)
		}

		app.Post(func() {
			ids := app.GetAvailableServerIDs(sgs_herox.SvrTypeGame)
			msg := &smsg.ReqGameRemove{}
			for _, id := range ids {
				cur := id
				app.GetServer(cur).ReqSugar(msg, func(resp *smsg.RespGameRemove, err error) {
					if err != nil {
						logrus.Error(fmt.Sprintf("ReqGameRemove:%d ", cur), err.Error())
						return
					}

					//logrus.Debug(fmt.Sprintf("ReqGameRemove:%d %d ", cur, resp.ErrCode), resp.GameIDs)
				}, 60*time.Second)
			}
		})
	})

	bs, _ := MakeCommonResult(0, "已提交审核", nil)
	w.Write(bs)
}

//拉玩家上内存（用于全服回档或部分玩家数据修改后，更新Redis用）
func (p *httpServer) RecoverRedis(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	defer p.write(res, w)

	logrus.Debug("RecoverRedis")
	useridstr := r.FormValue("userids")
	userids, err := gameutil.ParseUint64s(useridstr, ",")
	if err != nil {
		res["code"] = -1
		res["message"] = "参数错误"
		return
	}
	succuserids := make([]uint64, 0)

	//并行执行每次100个
	workerChan := make(chan struct{}, 100)
	failuseridchan := make(chan uint64, len(userids))
	wg := sync.WaitGroup{}
	wg.Add(len(userids))
	for _, userid := range userids {
		tmp := userid
		workerChan <- struct{}{}
		util.SafeGo(func() {
			defer func() {
				wg.Done()
				<-workerChan
			}()
			resp, err := entityMgr.RequestCall(&smsg.ReqUserSyncDB{
				Userid:   tmp,
				SyncType: 2,
			}, time.Second*10)

			if err != nil {
				failuseridchan <- tmp
				logrus.WithError(err).Error("RecoverRedis entityMgr.RequestCall")
				return
			}
			resp1 := resp.(*smsg.RespUserSyncDB)
			if resp1.ErrCode != 0 {
				failuseridchan <- tmp
				logrus.WithFields(logrus.Fields{
					"errcode": resp1.ErrCode,
					"userid":  tmp,
				}).WithError(err).Error("RecoverRedis resp1.ErrCode!=0")
				return
			}
		})

		succuserids = append(succuserids, userid)
	}
	wg.Wait()

	//收集错误信息
	failuserids := make([]uint64, 0)
	close(failuseridchan)
	for {
		failuserid, ok := <-failuseridchan
		if ok {
			failuserids = append(failuserids, failuserid)
		} else {
			break
		}
	}

	if len(failuserids) == 0 {
		res["code"] = 0
		res["message"] = "执行成功"
		return
	} else {
		res["code"] = 0
		res["message"] = "执行部分成功"
		res["failuserids"] = failuserids
		return
	}
}

func (p *httpServer) UserDataFrozen(w http.ResponseWriter, r *http.Request) {
	t, e := ParseHtmlFiles("web/template/user_data_frozen.html")
	if e != nil {
		fmt.Print("ParseFiles:", e)
		return
	}
	data := make(map[string]interface{})
	e = t.Execute(w, data)
	if e != nil {
		fmt.Print("Execute:", e)
		return
	}
}
func (p *httpServer) UserDataOpt(w http.ResponseWriter, r *http.Request) {
	t, e := ParseHtmlFiles("web/template/user_data_opt.html")
	if e != nil {
		fmt.Print("ParseFiles:", e)
		return
	}
	data := make(map[string]interface{})
	e = t.Execute(w, data)
	if e != nil {
		fmt.Print("Execute:", e)
		return
	}
}

func (p *httpServer) UserRedisDataOpt(w http.ResponseWriter, r *http.Request) {
	t, e := ParseHtmlFiles("web/template/user_redis_data_opt.html")
	if e != nil {
		fmt.Print("ParseFiles:", e)
		return
	}
	data := make(map[string]interface{})
	e = t.Execute(w, data)
	if e != nil {
		fmt.Print("Execute:", e)
		return
	}
}

func (p *httpServer) QueryUserInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	info := r.FormValue("info")
	info = strings.Trim(info, " ")

	if len(info) == 0 {
		w.Write([]byte("查询信息错误"))
		return
	}

	var l []string
	callback := func(s string) {
		if len(s) != 0 {
			l = append(l, s)
		}
	}
	gmMgr.onRemoteCommand("/q "+info, callback)

	tmp := strings.Join(l, "\n")
	w.Write([]byte(tmp))
}

func (p *httpServer) DoRemoteCommand(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	info := r.FormValue("info")
	info = strings.Trim(info, " ")

	if len(info) == 0 {
		bs, _ := MakeCommonResult(0, "时间错误", nil)
		w.Write(bs)
		return
	}
	var l []string
	callback := func(s string) {
		if len(s) != 0 {
			l = append(l, s)
		}
	}
	gmMgr.onRemoteCommand(info, callback)

	tmp := strings.Join(l, " ")
	bs, _ := MakeCommonResult(0, tmp, nil)
	w.Write([]byte(bs))
}

func (p *httpServer) FrozenUserData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	info := r.FormValue("info")
	info = strings.Trim(info, " ")

	userid, _ := strconv.ParseInt(info, 10, 64)
	logrus.Debug("FrozenUserData:", userid)
	if userid == 0 {
		bs, _ := MakeCommonResult(1, "没有找到角色", nil)
		w.Write(bs)
		return
	}
	dbMgr.QueryRawUserData(uint64(userid), "`userid`, `is_init`", func(rows *sql.Rows, err error) {
		if err != nil {
			bs, _ := MakeCommonResult(1, err.Error(), nil)
			w.Write(bs)
			return
		}
		if !rows.Next() {
			bs, _ := MakeCommonResult(1, "没有找到角色", nil)
			w.Write(bs)
			return
		}

		var uid uint64
		var is_init int32
		e := rows.Scan(&uid, &is_init)
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
		if is_init != 1 {
			bs, _ := MakeCommonResult(1, fmt.Sprintf("角色不可冻结 is_init:%d", is_init), nil)
			w.Write(bs)
			return
		}
		if is_init == 2 {
			bs, _ := MakeCommonResult(1, "角色已冻结", nil)
			w.Write(bs)
			return
		}

		////kick1
		//entityMgr.SendMsg(&smsg.AllLoReqKickUserOut{
		//	Userid: uint64(userid),
		//})
		//time.Sleep(5 * time.Second)

		//save
		resp, e := entityMgr.RequestCall(&smsg.ReqUserSyncDB{
			Userid:   uid,
			SyncType: 1,
		}, 15*time.Second)
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
		msg, _ := resp.(*smsg.RespUserSyncDB)
		if msg.ErrCode != 0 {
			bs, _ := MakeCommonResult(1, fmt.Sprintf("SyncDB Error:%d", msg.ErrCode), nil)
			w.Write(bs)
			return
		}
		//frozen
		tableName := gameshared.GetUserTableName(uint64(userid))
		sqlStr := fmt.Sprintf("update %s set `is_init` = 2 where userid = ?;", tableName)
		_, e = dbMgr.db.Exec(sqlStr, uint64(userid))
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
		//kick2
		//entityMgr.SendMsg(&smsg.AllLoReqKickUserOut{
		//	Userid: uint64(userid),
		//})
		//is_init被改变了 recover
		{
			_, e := entityMgr.RequestCall(&smsg.ReqUserSyncDB{
				Userid:   uid,
				SyncType: 2,
			}, 15*time.Second)
			if e != nil {
			}
		}
		bs, _ := MakeCommonResult(0, "ok", nil)
		w.Write(bs)
	})
}

func (p *httpServer) UnFrozenUserData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	info := r.FormValue("info")
	info = strings.Trim(info, " ")

	userid, _ := strconv.ParseInt(info, 10, 64)
	logrus.Debug("UnFrozenUserData:", userid)
	if userid == 0 {
		bs, _ := MakeCommonResult(1, "没有找到角色", nil)
		w.Write(bs)
		return
	}
	dbMgr.QueryRawUserData(uint64(userid), "`userid`, `is_init`", func(rows *sql.Rows, err error) {
		if err != nil {
			bs, _ := MakeCommonResult(1, err.Error(), nil)
			w.Write(bs)
			return
		}
		if !rows.Next() {
			bs, _ := MakeCommonResult(1, "没有找到角色", nil)
			w.Write(bs)
			return
		}

		var uid uint64
		var is_init int32
		e := rows.Scan(&uid, &is_init)
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
		if is_init != 2 {
			bs, _ := MakeCommonResult(1, fmt.Sprintf("角色不可解冻 is_init:%d", is_init), nil)
			w.Write(bs)
			return
		}

		//kick1
		//entityMgr.SendMsg(&smsg.AllLoReqKickUserOut{
		//	Userid: uint64(userid),
		//})
		//time.Sleep(5 * time.Second)
		//unfrozen
		tableName := gameshared.GetUserTableName(uint64(userid))
		sqlStr := fmt.Sprintf("update %s set `is_init` = 1 where userid = ? ;", tableName)
		_, e = dbMgr.db.Exec(sqlStr, uint64(userid))
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
		//recover
		resp, e := entityMgr.RequestCall(&smsg.ReqUserSyncDB{
			Userid:   uid,
			SyncType: 2,
		}, 15*time.Second)
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
		msg, _ := resp.(*smsg.RespUserSyncDB)
		if msg.ErrCode != 0 {
			bs, _ := MakeCommonResult(1, fmt.Sprint("SyncDB Error:%d", msg.ErrCode), nil)
			w.Write(bs)
			return
		}
		bs, _ := MakeCommonResult(0, "ok", nil)
		w.Write(bs)
	})
}

func (p *httpServer) RecoverUserData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	info := r.FormValue("info")
	info = strings.Trim(info, " ")

	userid, _ := strconv.ParseInt(info, 10, 64)
	logrus.Debug("RecoverUserData:", userid)
	if userid == 0 {
		bs, _ := MakeCommonResult(1, "没有找到角色", nil)
		w.Write(bs)
		return
	}
	dbMgr.QueryRawUserData(uint64(userid), "`userid`, `is_init`", func(rows *sql.Rows, err error) {
		if err != nil {
			bs, _ := MakeCommonResult(1, err.Error(), nil)
			w.Write(bs)
			return
		}
		if !rows.Next() {
			bs, _ := MakeCommonResult(1, "没有找到角色", nil)
			w.Write(bs)
			return
		}

		var uid uint64
		var is_init int32
		e := rows.Scan(&uid, &is_init)
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
		if is_init != 2 {
			bs, _ := MakeCommonResult(1, fmt.Sprint("角色未冻结 is_init:%d", is_init), nil)
			w.Write(bs)
			return
		}

		//kick1
		//entityMgr.SendMsg(&smsg.AllLoReqKickUserOut{
		//	Userid: uint64(userid),
		//})
		//time.Sleep(5 * time.Second)

		//recover
		resp, e := entityMgr.RequestCall(&smsg.ReqUserSyncDB{
			Userid:   uid,
			SyncType: 2,
		}, 15*time.Second)
		if e != nil {
			bs, _ := MakeCommonResult(1, e.Error(), nil)
			w.Write(bs)
			return
		}
		msg, _ := resp.(*smsg.RespUserSyncDB)
		if msg.ErrCode != 0 {
			bs, _ := MakeCommonResult(1, fmt.Sprint("SyncDB Error:%d", msg.ErrCode), nil)
			w.Write(bs)
			return
		}

		bs, _ := MakeCommonResult(0, "ok", nil)
		w.Write(bs)
	})
}
func (p *httpServer) QueryRawUserData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()
	info := r.FormValue("info")
	info = strings.Trim(info, " ")
	userid, _ := strconv.ParseInt(info, 10, 64)
	logrus.Debug("QueryRawUserData:", userid)
	if userid == 0 {
		//w.Write([]byte("id错误"))
		p.writeObj(&CommonResult{1, "id错误", nil}, w)
		return
	}

	content, md5, err := dbMgr.QueryAllRawUserDataWithMd5(uint64(userid), true)
	if err != nil {
		p.writeObj(&CommonResult{1, err.Error(), nil}, w)
		return
	}
	d := make(map[string]interface{})
	d["content"] = content
	d["md5"] = md5
	d["ret"] = 0
	bs, e := json.Marshal(d)
	if e != nil {
		p.writeObj(&CommonResult{1, e.Error(), nil}, w)
		return
	}
	w.Write(bs)
	return
	//dbMgr.QueryRawUserData(uint64(userid), "*", func(rows *sql.Rows, err error) {
	//	if err != nil {
	//		//w.Write([]byte(err.Error()))
	//		p.writeObj(&CommonResult{1, err.Error(), nil}, w)
	//		return
	//	}
	//	columns, err := rows.ColumnTypes()
	//	if err != nil {
	//		//w.Write([]byte(err.Error()))
	//		p.writeObj(&CommonResult{1, err.Error(), nil}, w)
	//		return
	//	}
	//
	//	if !rows.Next() {
	//		//w.Write([]byte(""))
	//		p.writeObj(&CommonResult{1, "no rows", nil}, w)
	//		return
	//
	//	}
	//	values := make([]interface{}, len(columns))
	//	scanArgs := make([]interface{}, len(values))
	//	for i := range values {
	//		scanArgs[i] = &values[i]
	//	}
	//	err = rows.Scan(scanArgs...)
	//	if err != nil {
	//		//w.Write([]byte(err.Error()))
	//		p.writeObj(&CommonResult{1, err.Error(), nil}, w)
	//		return
	//	}
	//	var columnsStringList []string
	//	for i, pval := range values {
	//		tmp := DBColumnInfo(columns[i], pval)
	//		columnsStringList = append(columnsStringList, tmp)
	//	}
	//	content := strings.Join(columnsStringList, "\n")
	//
	//	fpath := fmt.Sprintf("./raw/%s_%d.txt", time.Now().Format("20060102150405"), userid)
	//	err = os.MkdirAll(path.Dir(fpath), 0766)
	//	if err != nil {
	//		logrus.Error("QueryRawUserData Save " + fpath + " " + err.Error())
	//	} else {
	//		f, _ := os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC, 0666) //打开文件
	//		if f != nil {
	//			io.WriteString(f, content)
	//			f.Close()
	//		}
	//	}
	//	d := make(map[string]interface{})
	//	d["content"] = content
	//	d["md5"] = gameutil.MD5(content)
	//	d["ret"] = 0
	//	bs, e := json.Marshal(d)
	//	if e != nil {
	//		p.writeObj(&CommonResult{1, e.Error(), nil}, w)
	//		return
	//	}
	//	w.Write(bs)
	//	return
	//})
}

//操作玩家redis数据
func (p *httpServer) RedisDataOpt(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	defer p.write(res, w)
	r.ParseForm()
	//opt := r.FormValue("opt")
	key := r.FormValue("key")
	//content := r.FormValue("content")
	//md5 := r.FormValue("md5")
	//userID, err := gameutil.ToUInt64(r.FormValue("userid"))
	//if err != nil {
	//	res["code"] = -1
	//	res["message"] = "参数错误"
	//	return
	//}

	switch key {
	//case "friend":
	//	req := &smsg.AdFrReqFriendOption{
	//		TargetUserid: userID,
	//		Content:      content,
	//		LastMd5:      md5,
	//	}
	//	switch opt {
	//	case "query":
	//		req.Opt = smsg.AdFrReqFriendOption_Query
	//	case "save":
	//		req.Opt = smsg.AdFrReqFriendOption_Save
	//	default:
	//		res["code"] = -1
	//		res["message"] = "参数错误"
	//	}
	//	resp, err := app.GetService(sgs_herox.SvrTypeFriend).RequestCall(req, time.Second*30)
	//	if err != nil {
	//		logrus.WithField("userid", userID).WithError(err).Error("RedisDataOpt friend")
	//		res["code"] = -1
	//		res["message"] = "请求出错"
	//		p.write(res, w)
	//		return
	//	}
	//	resp1 := resp.(*smsg.AdFrRespFriendData)
	//	if resp1.ErrCode != 0 {
	//		logrus.WithField("userid", userID).WithError(err).Error("RedisDataOpt friend")
	//		res["code"] = -1
	//		res["message"] = resp1.ErrMsg
	//		return
	//	}
	//	switch req.Opt {
	//	case smsg.AdFrReqFriendOption_Query:
	//		res["content"] = resp1.Content
	//		res["md5"] = resp1.Md5
	//		res["userid"] = userID
	//		res["key"] = key
	//		return
	//	case smsg.AdFrReqFriendOption_Save:
	//		res["code"] = 0
	//		res["message"] = "请求执行成功"
	//	}
	default:
		res["code"] = -1
		res["message"] = "参数错误"
	}

	return
}
func (p *httpServer) SetColumnNull(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()
	info := r.FormValue("info")
	info = strings.Trim(info, " ")
	userid, _ := strconv.ParseInt(info, 10, 64)
	if userid == 0 {
		w.Write([]byte("设置column失败：id错误"))
		return
	}
	column := r.FormValue("column")
	column = strings.Trim(column, " ")
	if column == "" {
		w.Write([]byte("设置column失败：column错误"))
		return
	}

	md5 := r.FormValue("md5")
	if len(md5) == 0 {
		w.Write([]byte("设置column失败：md5错误"))
		return
	}
	_, oldMd5, err := dbMgr.QueryAllRawUserDataWithMd5(uint64(userid), false)
	if err != nil {
		p.writeObj(&CommonResult{1, err.Error(), nil}, w)
		w.Write([]byte("设置column失败：" + err.Error()))
		return
	}
	if md5 != oldMd5 {
		w.Write([]byte("设置column失败：md5错误"))
		return
	}

	dbMgr.QueryRawUserData(uint64(userid), "*", func(rows *sql.Rows, err error) {
		if err != nil {
			w.Write([]byte("设置column失败：" + err.Error()))
			return
		}
		if !rows.Next() {
			w.Write([]byte(fmt.Sprintf("设置column失败：未找到角色：%d", userid)))
			return
		}
		columns, err := rows.ColumnTypes()
		if err != nil {
			w.Write([]byte("设置column失败：" + err.Error()))
			return
		}

		var t *sql.ColumnType
		var idx int
		for i, c := range columns {
			if c.Name() == column {
				idx = i
				t = c
				break
			}
		}

		if t == nil {
			w.Write([]byte("设置column失败：未找到column：" + column))
			return
		}
		if a, b := t.Nullable(); !a || !b {
			w.Write([]byte(fmt.Sprintf("设置column失败：字段不支持NULL %s %s", column, t.ScanType().Name())))
			return
		}

		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)
		if err != nil {
			w.Write([]byte("设置column失败：" + err.Error()))
			return
		}

		//if t.ScanType().Name() != "NullTime" {
		//	w.Write([]byte(fmt.Sprintf("设置column失败：仅支持NullTime字段设置为NULL %s %s", column, t.ScanType().Name())))
		//	return
		//}

		tableName := gameshared.GetUserTableName(uint64(userid))
		sqlStr := fmt.Sprintf("update %s set `%s` = null where userid = ?;", tableName, column)
		_, e := dbMgr.db.Exec(sqlStr, uint64(userid))
		if e != nil {
			w.Write([]byte("设置column失败：" + e.Error()))
			return
		}

		tmp := DBColumnInfo(t, values[idx])
		logrus.Warnf("SetColumnNull: %d %s oldValue:%s", userid, column, tmp)

		w.Write([]byte(fmt.Sprintf("设置column完成：%s为NULL", column)))
		return
	})
}
func (p *httpServer) SetColumnData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()
	info := r.FormValue("info")
	info = strings.Trim(info, " ")
	userid, _ := strconv.ParseInt(info, 10, 64)
	if userid == 0 {
		w.Write([]byte("设置column失败：id错误"))
		return
	}

	column := r.FormValue("column")
	column = strings.Trim(column, " ")
	if column == "" {
		w.Write([]byte("设置column失败：column错误"))
		return
	}

	md5 := r.FormValue("md5")
	if len(md5) == 0 {
		w.Write([]byte("设置column失败：md5错误"))
		return
	}
	_, oldMd5, err := dbMgr.QueryAllRawUserDataWithMd5(uint64(userid), false)
	if err != nil {
		p.writeObj(&CommonResult{1, err.Error(), nil}, w)
		w.Write([]byte("设置column失败：" + err.Error()))
		return
	}
	if md5 != oldMd5 {
		w.Write([]byte("设置column失败：md5错误"))
		return
	}

	columnData := r.FormValue("column_data")
	//logrus.Debugf("SetColumnData: %d %s %s", userid, column, columnData)
	dbMgr.QueryRawUserData(uint64(userid), "*", func(rows *sql.Rows, err error) {
		if err != nil {
			w.Write([]byte("设置column失败：" + err.Error()))
			return
		}
		if !rows.Next() {
			w.Write([]byte(fmt.Sprintf("设置column失败：未找到角色：%d", userid)))
			return
		}
		columns, err := rows.ColumnTypes()
		if err != nil {
			w.Write([]byte("设置column失败：" + err.Error()))
			return
		}
		var t *sql.ColumnType
		var idx int
		for i, c := range columns {
			if c.Name() == column {
				idx = i
				t = c
				break
			}
		}
		if t == nil {
			w.Write([]byte("设置column失败：未找到column：" + column))
			return
		}

		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)
		if err != nil {
			w.Write([]byte("设置column失败：" + err.Error()))
			return
		}

		tableName := gameshared.GetUserTableName(uint64(userid))
		switch t.ScanType().Name() {
		case "NullTime":
			if len(columnData) == 0 {
				w.Write([]byte("设置column失败：" + "NullTime字段内容为空"))
				return
			}
			_, e := time.ParseInLocation("2006-01-02 15:04:05", columnData, time.Local)
			if e != nil {
				w.Write([]byte("设置column失败：" + e.Error()))
				return
			}
			sqlStr := fmt.Sprintf("update %s set `%s` = ? where userid = ?;", tableName, column)
			_, e = dbMgr.db.Exec(sqlStr, columnData, uint64(userid))
			if e != nil {
				w.Write([]byte("设置column失败：" + e.Error()))
				return
			}
		case "RawBytes":
			if len(columnData) == 0 {
				sqlStr := fmt.Sprintf("update %s set `%s` = '' where userid = ?;", tableName, column)
				_, e := dbMgr.db.Exec(sqlStr, uint64(userid))
				if e != nil {
					w.Write([]byte("设置column失败：" + e.Error()))
					return
				}
			} else {
				bs, e := hex.DecodeString(columnData)
				if e != nil {
					w.Write([]byte("设置column失败：" + e.Error()))
					return
				}
				sqlStr := fmt.Sprintf("update %s set `%s` = ? where userid = ?;", tableName, column)
				_, e = dbMgr.db.Exec(sqlStr, bs, uint64(userid))
				if e != nil {
					w.Write([]byte("设置column失败：" + e.Error()))
					return
				}
			}
		case "int8", "int16", "int32", "int64", "int", "uint8", "uint16", "uint32", "uint64", "uint":
			var v int64
			var e error
			if len(columnData) != 0 {
				v, e = strconv.ParseInt(columnData, 10, 64)
				if e != nil {
					w.Write([]byte("设置column失败：" + e.Error()))
					return
				}
			}
			sqlStr := fmt.Sprintf("update %s set `%s` = ? where userid = ?;", tableName, column)
			_, e = dbMgr.db.Exec(sqlStr, v, uint64(userid))
			if e != nil {
				w.Write([]byte("设置column失败：" + e.Error()))
				return
			}
		default:
			w.Write([]byte(fmt.Sprintf("设置column失败：该字段类型尚未支持 %s %s", column, t.ScanType().Name())))
			return
		}
		tmp := DBColumnInfo(t, values[idx])
		logrus.Warnf("SetColumnData: %d %s oldValue:%s newData:%s", userid, column, tmp, columnData)
		//logrus.Debugf("SetColumnData: %d %s %s", userid, column, columnData)
		w.Write([]byte(fmt.Sprintf("设置column完成：%s", column)))
		return
	})
}
func (p *httpServer) QueryRawToContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()

	column := r.FormValue("column")
	column = strings.Trim(column, " ")
	data := r.FormValue("data")
	if len(data) == 0 {
		w.Write([]byte("转换失败：空数据"))
		return
	}
	if column == "" {
		tmp, e := hex.DecodeString(data)
		if e != nil {
			w.Write([]byte("转换失败：未指定column,按字符串转码失败 " + e.Error()))
			return
		}
		w.Write([]byte(string(tmp)))
		return
	}

	dbMgr.QueryRawUserData(gameshared.UserIDBase, "*", func(rows *sql.Rows, err error) {
		if err != nil {
			w.Write([]byte("转换失败：" + err.Error()))
			return
		}
		columns, err := rows.ColumnTypes()
		if err != nil {
			w.Write([]byte("转换失败：" + err.Error()))
			return
		}
		var t *sql.ColumnType
		for _, c := range columns {
			if c.Name() == column {
				t = c
				break
			}
		}
		if t == nil {
			w.Write([]byte("转换失败：未找到column：" + column))
			return
		}
		bs, e := hex.DecodeString(data)
		if e != nil {
			w.Write([]byte("转换失败：" + e.Error()))
			return
		}
		switch t.ScanType().Name() {
		case "NullTime":
			w.Write([]byte("转换失败：该类型数据不需要转换" + column + " " + t.ScanType().Name()))
		case "RawBytes":
			member := gameutil.NewStructMember(&userdb.User{}, column)
			if member != nil {
				switch member.(type) {
				case string, *string:
					w.Write(bs)
					return
				default:
					tmp, e := ProtoToJson(bs, member.(proto.Message))
					if e == nil {
						w.Write(tmp)
						return
					}
				}
			}
			content := ColumnProtoMessage(column)
			if content != nil {
				tmp, e := ProtoToJson(bs, content)
				if e != nil {
					w.Write([]byte("转换失败：" + e.Error()))
					return
				}
				w.Write(tmp)
				return
			}
			w.Write([]byte("转换失败：该字段类型尚未支持" + column + " " + t.ScanType().Name()))
			return
		case "int8", "int16", "int32", "int64", "int", "uint8", "uint16", "uint32", "uint64", "uint":
			w.Write([]byte("转换失败：该类型数据不需要转换" + column + " " + t.ScanType().Name()))
		default:
			w.Write([]byte("转换失败：该字段类型尚未支持" + column + " " + t.ScanType().Name()))
			return
		}
	})
}
func (p *httpServer) QueryContentToRaw(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()

	column := r.FormValue("column")
	column = strings.Trim(column, " ")
	data := r.FormValue("data")
	if len(data) == 0 {
		w.Write([]byte("转换失败：空数据"))
		return
	}
	if column == "" {
		tmp := []byte(data)
		s_txt := fmt.Sprintf("%x", tmp)
		w.Write([]byte(s_txt))
		return
	}

	dbMgr.QueryRawUserData(gameshared.UserIDBase, "*", func(rows *sql.Rows, err error) {
		if err != nil {
			w.Write([]byte("转换失败：" + err.Error()))
			return
		}
		columns, err := rows.ColumnTypes()
		if err != nil {
			w.Write([]byte("转换失败：" + err.Error()))
			return
		}
		var t *sql.ColumnType
		for _, c := range columns {
			if c.Name() == column {
				t = c
				break
			}
		}
		if t == nil {
			w.Write([]byte("转换失败：未找到column：" + column))
			return
		}
		switch t.ScanType().Name() {
		case "NullTime":
			w.Write([]byte("转换失败：该类型数据不需要转换" + column + " " + t.ScanType().Name()))
		case "RawBytes":
			member := gameutil.NewStructMember(&userdb.User{}, column)
			if member != nil {
				switch member.(type) {
				case string, *string:
					s_txt := fmt.Sprintf("%x", data)
					w.Write([]byte(s_txt))
					return
				default:
					tmp, e := JsonToProto(data, member.(proto.Message))
					if e == nil {
						w.Write([]byte(tmp))
						return
					}
				}
			}
			content := ColumnProtoMessage(column)
			if content != nil {
				tmp, e := JsonToProto(data, content)
				if e != nil {
					w.Write([]byte("转换失败：" + e.Error()))
					return
				}
				w.Write([]byte(tmp))
				return
			}
			w.Write([]byte("转换失败：该字段类型尚未支持" + column + " " + t.ScanType().Name()))
			return
		case "int8", "int16", "int32", "int64", "int", "uint8", "uint16", "uint32", "uint64", "uint":
			w.Write([]byte("转换失败：该类型数据不需要转换" + column + " " + t.ScanType().Name()))
		default:
			w.Write([]byte("转换失败：该字段类型尚未支持" + column + " " + t.ScanType().Name()))
			return
		}
	})
}
func (p *httpServer) QueryZipRawToContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()

	column := r.FormValue("column")
	column = strings.Trim(column, " ")
	data := r.FormValue("data")
	if len(data) == 0 {
		w.Write([]byte("转换失败：空数据"))
		return
	}
	if column != "" {
		w.Write([]byte("转换失败：Base64Unzip转码不能指定列"))
		return
	}
	//tmp, e := hex.DecodeString(data)
	//if e != nil {
	//	w.Write([]byte("转换失败：hex转码失败 " + e.Error()))
	//	return
	//}
	bs, e := gameutil.UnZipBase64(data)
	if e != nil {
		w.Write([]byte("转换失败：Base64Unzip转码失败 " + e.Error()))
		return
	}
	w.Write([]byte(string(bs)))
	return
}
func (p *httpServer) QueryContentToZipRaw(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()

	column := r.FormValue("column")
	column = strings.Trim(column, " ")
	data := r.FormValue("data")
	if len(data) == 0 {
		w.Write([]byte("转换失败：空数据"))
		return
	}
	if column != "" {
		w.Write([]byte("转换失败：ZipBase64转码不能指定列"))
		return
	}
	bs, e := gameutil.ZipBase64([]byte(data))
	if e != nil {
		w.Write([]byte("转换失败：ZipBase64转码失败 " + e.Error()))
		return
	}
	//tmp := hex.EncodeToString(bs)
	w.Write([]byte(bs))
	return
}
func ProtoToJson(bs []byte, v proto.Message) ([]byte, error) {
	e := proto.Unmarshal(bs, v)
	if e != nil {
		return nil, e
	}
	tmp, e := json.MarshalIndent(v, "", "\t")
	if e != nil {
		return nil, e
	}
	return tmp, nil
}
func JsonToProto(data string, v proto.Message) (string, error) {
	e := json.Unmarshal([]byte(data), v)
	if e != nil {
		return "", e
	}
	tmp, e := proto.Marshal(v)
	if e != nil {
		return "", e
	}
	s_txt := fmt.Sprintf("%x", tmp)
	return s_txt, nil
}
func (p *httpServer) UserDataEnc(w http.ResponseWriter, r *http.Request) {
	t, e := ParseHtmlFiles("web/template/user_data_enc.html")
	if e != nil {
		fmt.Print("ParseFiles:", e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Print("Execute:", e)
		return
	}
}

func (p *httpServer) ServerInfo(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	defer p.write(res, w)

	infoMap, err := dbMgr.GetServerInfo()
	if err != nil {
		res["error"] = err.Error()
		return
	}

	r.ParseForm()
	key := r.FormValue("key")
	key = strings.Trim(key, " ")
	value := r.FormValue("value")
	value = strings.Trim(value, " ")

	dirty, old, err := CheckServerInfo(infoMap, key, value)
	if err != nil {
		res["error"] = err.Error()
		return
	}
	if !dirty {
		res["Sync"] = true
	}
	if infoMap != nil {
		res["Infos"] = infoMap
	}
	if old != value {
		res["change"] = fmt.Sprintf("[%v]%v->%v", key, old, value)
	}

	if key == "DingDingTokens" {
		var tokens []string
		if len(strings.TrimSpace(value)) > 0 {
			tokens = strings.Split(value, ";")
		}
		reporterMgr, err = newReporterManager(tokens)
		if err != nil {
			logrus.WithError(err).Errorf("newreporterManger %s", tokens)
		}
	}
}
