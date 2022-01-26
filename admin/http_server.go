package admin

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type httpServer struct {
	listenAddress string
	serverMux     *http.ServeMux
	server        *http.Server
	authKey       string
}

func (p *httpServer) init(listenAddr string) {
	p.listenAddress = listenAddr
	p.serverMux = http.NewServeMux()

	//promethous
	p.serverMux.Handle("/metrics", promhttp.Handler())

	//----------------------------------------------对公共组接口-------------------------------------------------------
	// 后面所有的GM命令都 添加到这里, 便于聊天和后台同时可用
	p.serverMux.HandleFunc("/v1/gm_command", p.authOnly(p.gmCommand))

	//提供给运维组重载配置的接口
	p.serverMux.HandleFunc("/v2/reload", p.authOnly(ReloadV2))

	files := http.FileServer(http.Dir("web"))
	p.serverMux.Handle("/js/", files)
	p.serverMux.Handle("/depend/", files)
	p.serverMux.Handle("/css/", files)
	p.serverMux.Handle("/img/", files)

	p.serverMux.HandleFunc("/", HttpPlat)
	p.serverMux.HandleFunc("/home_page", HomePage)
	p.serverMux.HandleFunc("/Register", Register)
	p.serverMux.HandleFunc("/login", Login)
	p.serverMux.HandleFunc("/logout", p.loginOnly(Logout))
	p.serverMux.HandleFunc("/password", p.loginOnly(Password))
	p.serverMux.HandleFunc("/password_edit", p.loginOnlyWithAccount(PasswordEdit))

	//角色管理
	p.registerAdminHandler("/user_data_frozen", "角色：管理", 10100, 1, p.UserDataFrozen)
	p.registerAdminHandler("/v9/FrozenUserData", "冻结角色", 10101, 0, p.FrozenUserData)
	p.registerAdminHandler("/v9/UnFrozenUserData", "解冻角色", 10102, 0, p.UnFrozenUserData)
	p.registerAdminHandler("/v9/RecoverUserData", "重载角色数据", 10103, 0, p.RecoverUserData)
	p.registerAdminHandler("/v0/QueryUserInfo", "角色：查询", 10104, 0, p.QueryUserInfo)

	//TODO 删除多余的
	p.registerAdminHandler("/v9/DoRemoteCommand", "远程GM", 10123, 0, p.DoRemoteCommand)

	p.registerAdminHandler("/user_data_opt", "角色：数据操作", 10200, 1, p.UserDataOpt)
	p.registerAdminHandler("/v9/SetColumnNull", "设置字段为NULL", 10201, 0, p.SetColumnNull)
	p.registerAdminHandler("/v9/SetColumnData", "设置数据到字段", 10202, 0, p.SetColumnData)
	p.registerAdminHandler("/v9/QueryRawUserData", "数据操作:查询角色信息", 10203, 0, p.QueryRawUserData)

	p.registerAdminHandler("/user_data_enc", "角色：数据解析", 10300, 1, p.UserDataEnc)
	p.registerAdminHandler("/v9/QueryRawToContent", "Raw->内容", 10301, 0, p.QueryRawToContent)
	p.registerAdminHandler("/v9/QueryContentToRaw", "内容->Raw", 10302, 0, p.QueryContentToRaw)
	p.registerAdminHandler("/v9/QueryZipRawToContent", "Base64Unzip", 10303, 0, p.QueryZipRawToContent)
	p.registerAdminHandler("/v9/QueryContentToZipRaw", "ZipBase64", 10304, 0, p.QueryContentToZipRaw)

	p.registerAdminHandler("/user_redis_data_opt", "角色：Redis数据操作", 10400, 1, p.UserRedisDataOpt)
	p.registerAdminHandler("/v9/RecoverRedisData", "刷新玩家Reids", 10401, 0, p.RecoverRedis)
	p.registerAdminHandler("/v9/RedisDataOpt", "操作玩家Redis数据", 10402, 0, p.RedisDataOpt)

	//服务管理
	p.registerAdminHandler("/server_ctrl", "服务器：整体控制", 20100, 1, p.ServerCtrl)
	p.registerAdminHandler("/v9/KickAll", "所有人下线", 20101, 0, p.KickAll)
	p.registerAdminHandler("/v9/Silence", "所有人下线，结束所有游戏", 20102, 0, p.KickAllAndForceGameOver)
	p.registerAdminHandler("/v9/CloseServer", "关闭服务器", 20103, 0, p.closeServer)

	p.registerAdminHandler("/reloadpage", "服务器：重载配置", 20200, 1, ReloadPage)
	p.registerAdminHandler("/reload", "重载配置", 20201, 0, Reload)
	p.registerAdminHandler("/server_config", "服务器：设置", 20300, 1, p.ServerConfig)
	p.registerAdminHandler("/v2/serverInfo", "服务器：配置", 20301, 0, p.ServerInfo)

	p.registerAdminHandler("/v9/EditServerOpenTime", "修改开服时间", 20302, 0, p.EditServerOpenTime)
	p.registerAdminHandler("/v9/EditServerCloseTime", "修改关服时间", 20303, 0, p.EditServerCloseTime)
	p.registerAdminHandler("/v9/EditInnerIp", "修改内网ip", 20304, 0, p.EditInnerIp)

	//审核
	p.registerAdminHandler("/action_verify", "审核", 70100, 1, p.ActionVerify)
	p.registerAdminHandler("/v9/DoActionVerify", "审核", 70101, 0, p.DoActionVerify)

	//后台账号管理
	p.registerAdminHandlerWithLoginInfo("/admin_set", "admin账号管理", 80100, 1, admin_set)
	p.registerAdminHandler("/AddGM", "添加账号", 80101, 0, AddGM)
	p.registerAdminHandler("/UpdateGM", "修改账号", 80102, 0, UpdateGM)
	p.registerAdminHandler("/RemoveGM", "删除账号", 80103, 0, RemoveGM)

	//其它
	p.registerAdminHandlerWithLoginInfo("/CMD", "控制台", 90100, 1, CMD)

	p.registerAdminHandlerWithLoginInfo("/MsgMonitorPage", "监控控制台", 90300, 1, MsgMonitorPage)
	p.registerAdminHandler("/wsmsgmonitor", "监控控制台", 90301, 0, WsMsgMonitorConnect)
}

func (p *httpServer) registerAdminHandler(pattern string, name string, rightID int, view int, handler func(http.ResponseWriter, *http.Request)) {
	right := &Rights{rightID, name, pattern, view}
	AllRights = append(AllRights, right)
	p.serverMux.HandleFunc(pattern, p.adminOnly(handler))
}

func (p *httpServer) registerAdminHandlerWithLoginInfo(pattern string, name string, rightID int, view int, handler func(http.ResponseWriter, *http.Request, *LoginInfo)) {
	right := &Rights{rightID, name, pattern, view}
	AllRights = append(AllRights, right)
	p.serverMux.HandleFunc(pattern, p.adminOnlyWithLoginInfo(handler))
}

//由于nginx默认过滤 header中带下划线的头，这里定义新的Auth头
func (p *httpServer) authOnly(h http.HandlerFunc) http.HandlerFunc {
	type response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//开发期暂不启用权限校验
		authKey := r.Header.Get("Auth-Key")
		if authKey != p.authKey {
			p.respond(w, r, response{Code: -1, Msg: "权限不足"}, http.StatusOK)
			return
		}
		h(w, r)
		return
	}
}

func (p *httpServer) loginOnly(h http.HandlerFunc) http.HandlerFunc {
	type response struct {
		Code   int    `json:"code"`
		Reason string `json:"Reason"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//s := sessionMgr.SessionStart(w, r)
		//account := s.Get("account")
		//if account == nil {
		//	p.respond(w, r, response{Code: -1, Reason: "登录失败，请重新登录"}, http.StatusOK)
		//	return
		//}
		cacheInfo,_ :=ioutil.ReadFile(SessionRoot)
		var mp map[string](time.Time)
		json.Unmarshal(cacheInfo,&mp)
		if len(mp)==0{
				p.respond(w, r, response{Code: -1, Reason: "登录失败，请重新登录"}, http.StatusOK)
				return
		}


		h(w, r)
		return
	}
}

func (p *httpServer) loginOnlyWithAccount(h func(w http.ResponseWriter, r *http.Request, account string)) http.HandlerFunc {
	type response struct {
		Code   int    `json:"code"`
		Reason string `json:"Reason"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//s := sessionMgr.SessionStart(w, r)
		//account := s.Get("account")
		//if account == nil {
		//	p.respond(w, r, response{Code: -1, Reason: "登录失败，请重新登录"}, http.StatusOK)
		//	return
		//}
		cacheInfo,_ :=ioutil.ReadFile(SessionRoot)
		var mp map[string](time.Time)
		json.Unmarshal(cacheInfo,&mp)
		if len(mp)==0{
			p.respond(w, r, response{Code: -1, Reason: "登录失败，请重新登录"}, http.StatusOK)
			return
		}
		var key string

		for key,_=range mp{}
		h(w, r, key)
		return
	}
}

func (p *httpServer) adminOnlyWithLoginInfo(h func(http.ResponseWriter, *http.Request, *LoginInfo)) http.HandlerFunc {
	type response struct {
		Code   int    `json:"code"`
		Reason string `json:"Reason"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		pattern := r.RequestURI

		//s := sessionMgr.SessionStart(w, r)
		//account := s.Get("account")
		//if account == nil {
		//	p.respond(w, r, response{Code: -1, Reason: "登录失败，请重新登录"}, http.StatusOK)
		//	return
		//}
		cacheInfo,_ :=ioutil.ReadFile(SessionRoot)
		var mp map[string](time.Time)
		json.Unmarshal(cacheInfo,&mp)
		if len(mp)==0{
			p.respond(w, r, response{Code: -1, Reason: "登录失败，请重新登录"}, http.StatusOK)
			return
		}
		var key string
		for key,_=range mp{}
		ci := GetLoginInfoByAccount(key)
		if ci == nil || !ci.CheckRights(pattern) {
			p.respond(w, r, response{Code: -1, Reason: "权限不足"}, http.StatusOK)
			return
		}
		h(w, r, ci)
		return
	}
}

func (p *httpServer) adminOnly(h http.HandlerFunc) http.HandlerFunc {

	type response struct {
		Code   int    `json:"code"`
		Reason string `json:"Reason"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		pattern := r.RequestURI

		//s := sessionMgr.SessionStart(w, r)
		//account := s.Get("account")
		//if account == nil {
		//	p.respond(w, r, response{Code: -1, Reason: "登录失败，请重新登录"}, http.StatusOK)
		//	return
		//}
		//
		//ci := GetLoginInfoByAccount(account.(string))
		cacheInfo,_ :=ioutil.ReadFile(SessionRoot)
		var mp map[string](time.Time)
		json.Unmarshal(cacheInfo,&mp)
		if len(mp)==0{
			p.respond(w, r, response{Code: -1, Reason: "登录失败，请重新登录"}, http.StatusOK)
			return
		}
		var key string
		for key,_=range mp{}
		ci := GetLoginInfoByAccount(key)
		if ci == nil || !ci.CheckRights(pattern) {
			p.respond(w, r, response{Code: -1, Reason: "权限不足"}, http.StatusOK)
			return
		}
		h(w, r)
		return
	}
}

func (p *httpServer) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logrus.WithField("data", data).WithError(err).Error("Server respond json encode error")
		}
	}
}

func (p *httpServer) run() {
	//TODO 权限认证
	p.server = &http.Server{Addr: p.listenAddress, Handler: p.serverMux}
	err := p.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logrus.WithFields(logrus.Fields{
			"listenAddress": p.listenAddress,
		}).WithError(err).Error("http.ListenAndServer error")
		panic(err)
	}
}

func (p *httpServer) close() {
	p.server.Close()
}
