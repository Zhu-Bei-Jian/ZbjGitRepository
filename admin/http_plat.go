package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/proto/smsg"
	"sort"
	"sync"
	"time"
)

var connectID int64
var connectMap map[int64]string

func GetConnectMap() map[int64]string {
	if connectMap == nil {
		connectMap = make(map[int64]string)
	}
	return connectMap
}

type CommonResult struct {
	Ret    int
	Reason string
	Data   interface{}
}

type Response struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func MakeCommonResult(ret int, reason string, i interface{}) ([]byte, error) {
	out := &CommonResult{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func MakeJson(i interface{}) ([]byte, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func HttpWrite(res http.ResponseWriter, f string, d interface{}) {
	t, e := ParseHtmlFiles(f)
	if e != nil {
		fmt.Print("ParseFiles:", e)
		return
	}
	e = t.Execute(res, d)
	if e != nil {
		fmt.Print("Execute:", e)
		return
	}
}
func HttpWriteJson(res http.ResponseWriter, d interface{}) {
	b, err := json.Marshal(d)
	if err != nil {
		return
	}
	n, err := res.Write(b)
	fmt.Println(n, err)
	return
}
func HttpPlat(w http.ResponseWriter, r *http.Request) {
	//s := sessionMgr.SessionStart(w, r)
	//account := s.Get("account")
	//ci := GetLoginInfo(req)
	cacheInfo, _ := ioutil.ReadFile(SessionRoot)
	var mp map[string](time.Time)
	json.Unmarshal(cacheInfo, &mp)
	if len(mp) == 0 {
		HttpWrite(w, "web/template/login.html", nil)
		return
	}
	var (
		key string
		val time.Time
	)
	for key, val = range mp {
	}
	nowTime := time.Now()
	//dTime :=nowTime.Sub(val)
	//var maxLiftTime float64=3600
	if nowTime.Sub(val).Seconds() > maxLiftTime {
		delete(mp, key)
		jsonBytes, _ := json.Marshal(mp)
		ioutil.WriteFile(SessionRoot, jsonBytes, 0664)

		HttpWrite(w, "web/template/login.html", nil)
		return
	}
	//if account == nil {
	//	HttpWrite(w, "web/template/login.html", nil)
	//	return
	//}
	ci := GetLoginInfoByAccount(key)
	Data := make(map[string]interface{})
	Data["Area"] = fmt.Sprintf("%d服", AreaId)
	Data["Name"] = ci.Account
	Data["IsAdmin"] = ci.IsAdmin
	//Data["RightsGroup"] = "1"
	//Data["Rights"] = []*Rights{
	//	//
	//	&Rights{1, "控制台", "/CMD", 1},
	//
	//	&Rights{1001, "角色：管理", "/user_data_frozen", 1},
	//	&Rights{1002, "角色：数据操作", "/user_data_opt", 1},
	//	&Rights{1003, "角色：数据解析", "/user_data_enc", 1},
	//	&Rights{1004, "角色：Redis数据操作", "/user_redis_data_opt", 1},
	//	//配置类
	//	&Rights{2001, "重载配置", "/reloadpage", 1},
	//
	//	&Rights{2002, "服务器配置", "/server_config", 1},
	//	//&Rights{200201, "服务器配置同步", "/v2/serverInfo", 1},
	//	&Rights{2003, "服务器控制", "/server_ctrl", 1},
	//
	//	//测试
	//	&Rights{3001, "胜率测试", "/wintest", 1},
	//
	//	&Rights{4001, "审核", "/action_verify", 1},
	//	&Rights{5001, "设置", "/admin_set", 1},
	//}
	var l []*Rights
	for _, r := range ci.rights {
		l = append(l, r)
	}
	sort.Sort(SortRights{l})
	Data["Rights"] = l
	//Data["ServerGroup"] = "s"

	HttpWrite(w, "web/template/plat.html", Data)
}

func HomePage(res http.ResponseWriter, req *http.Request) {
	HttpWrite(res, "web/template/homepage.html", nil)
}

func Reload(response http.ResponseWriter, req *http.Request) {
	opt := req.PostFormValue("opt")
	DoReload(response, opt, false /**v2**/)
}

func ReloadV2(response http.ResponseWriter, req *http.Request) {
	opt := req.FormValue("opt")
	DoReload(response, opt, true)
}

func DoReload(response http.ResponseWriter, optType string, v2 bool) {
	res := make(map[string]interface{})
	serverTypes := make([]appframe.ServerType, 0)

	switch optType {
	case "all":
		serverTypes = reloadServers
	case "auth":
		serverTypes = append(serverTypes, sgs_herox.SvrTypeAuth)
	case "lobby":
		serverTypes = append(serverTypes, sgs_herox.SvrTypeLobby)
	case "game":
		serverTypes = append(serverTypes, sgs_herox.SvrTypeGame)
	case "entity":
		serverTypes = append(serverTypes, sgs_herox.SvrTypeEntity)
	default:
		res["错误"] = "参数错误"
		return
	}

	//request all server，记录request servers,并等待全部返回
	var wg sync.WaitGroup
	allReloadServerIds := make([]uint32, 0)
	succServerIds := make([]uint32, 0)
	failServerIds := make([]uint32, 0)
	for _, v := range serverTypes {
		serverIDs := app.GetAvailableServerIDs(v)
		for _, v := range serverIDs {
			wg.Add(1)
			serverID := v
			allReloadServerIds = append(allReloadServerIds, serverID)
			app.GetServer(serverID).ReqSugar(&smsg.AdAllReqReload{}, func(resp *smsg.AdAllRespReload, err error) {
				defer wg.Done()

				if err != nil {
					failServerIds = append(failServerIds, serverID)
					return
				}

				if resp.ErrCode != smsg.AdAllRespReload_Invalid {
					failServerIds = append(failServerIds, serverID)
					return
				}
				succServerIds = append(succServerIds, serverID)
			}, time.Second*60)
		}
	}

	fmt.Println("allReloadServerID:", allReloadServerIds)
	wg.Wait()
	//logrus.Debugf("reloadSuccessServerIDs:%v,reloadFailServerIDs:%v", succServerIds, failServerIds)

	ret := fmt.Sprintf("成功:%v 失败:%v", succServerIds, failServerIds)

	if v2 {
		result := Response{Code: 0, Message: ""}
		if len(failServerIds) > 0 {
			result.Code = 1
		}
		result.Message = fmt.Sprintf("成功:%v 失败:%v", succServerIds, failServerIds)
		HttpWriteJson(response, result)
	} else {
		HttpWriteJson(response, &CommonResult{1, ret, nil})
	}

	reporterMgr.Send(fmt.Sprintf("服务器【%v】重载配置操作结果,%s", AreaId, ret))
}

func ReloadPage(res http.ResponseWriter, req *http.Request) {
	//ci := GetLoginInfo(req)
	//if ci == nil {
	//	HttpWrite(res, "web/template/error.html", &CommonResult{0, "登录信息已过期", nil})
	//	return
	//}
	HttpWrite(res, "web/template/reloadpage.html", nil)
}

func CMD(res http.ResponseWriter, req *http.Request, l *LoginInfo) {
	//ci := GetLoginInfo(req)
	//if ci == nil || !ci.CheckRights("/CMD") {
	//	HttpWrite(res, "web/template/error.html", &CommonResult{0, "登录信息已过期", nil})
	//	return
	//}

	//connectID++
	//
	//l.ConnectID = connectID
	//GetConnectMap()[l.ConnectID] = l.Account
	//
	//t, e := ParseHtmlFiles("web/template/monitor.html")
	//if e != nil {
	//	fmt.Print("ParseFiles:", e)
	//	HttpWrite(res, "web/template/error.html", &CommonResult{0, e.Error(), nil})
	//	return
	//}
	//
	//connInfo := &ConnectInfo{Name: l.Account, CID: l.ConnectID}
	//b, e := json.Marshal(connInfo)
	//if e != nil {
	//	HttpWrite(res, "web/template/error.html", &CommonResult{0, e.Error(), nil})
	//	return
	//}
	//data := &WsArgMonitor{Addr: MAddr, ConnectInfo: string(b)}
	//e = t.Execute(res, data)
	//if e != nil {
	//	fmt.Print("Execute:", e)
	//	return
	//}
}
