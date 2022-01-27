package admin

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"sanguosha.com/sgs_herox/gameutil"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_ADMIN     = "a"
	GROUP_TAG     int = 10000
)

type Rights struct {
	Id       int
	Name     string
	Action   string
	ViewSort int //菜单分类 0不在菜单中显示
}
type SortRights struct {
	l []*Rights
}

func (ts SortRights) Len() int {
	return len(ts.l)
}
func (ts SortRights) Swap(i, j int) {
	ts.l[i], ts.l[j] = ts.l[j], ts.l[i]
}

func (ts SortRights) Less(i, j int) bool {
	return ts.l[i].Id < ts.l[j].Id
}

type GMInfo struct {
	Name        string
	Pwd         string
	Rights      string
	ServerGroup string
	UserType    int
}

type GMInfoEx struct {
	Rights []int
}

var AllRights []*Rights

//var AllRights = []*Rights{
//	//角色管理
//	&Rights{10100, "角色：管理", "/user_data_frozen", 1},
//	&Rights{10200, "角色：数据操作", "/user_data_opt", 1},
//	&Rights{10300, "角色：数据解析", "/user_data_enc", 1},
//	&Rights{10400, "角色：Redis数据操作", "/user_redis_data_opt", 1},
//
//	//服务管理
//	&Rights{20100, "服务器：整体控制", "/server_ctrl", 1},
//	&Rights{20200, "服务器：重载配置", "/reloadpage", 1},
//	&Rights{20300, "服务器：设置", "/server_config", 1},
//
//	//审核
//	&Rights{70100, "审核", "/action_verify", 1},
//
//	//设置
//	&Rights{80100, "admin账号管理", "/admin_set", 1},
//	&Rights{80101, "-添加账号", "/AddGM", 0},
//	&Rights{80102, "-修改账号", "/UpdateGM", 0},
//	&Rights{80103, "-删除账号", "/RemoveGM", 0},
//
//	//其他
//	&Rights{90100, "控制台", "/CMD", 1},
//	&Rights{90200, "胜率测试", "/WinTest", 1},
//}

type LoginInfo struct {
	//Name  string
	//Token string

	Account string

	rights map[int]*Rights

	ConnectID int64
	IsAdmin   int
}

func (li *LoginInfo) CheckRights(action string) bool {
	for _, v := range li.rights {
		if v.Action == action {
			return true
		}
	}
	return false
}

func GetLoginInfoByAccount(account string) *LoginInfo {
	if account == "" {
		return nil
	}

	li := &LoginInfo{Account: account}
	li.rights = make(map[int]*Rights)
	if account == DEFAULT_ADMIN {
		for _, r := range AllRights {
			li.rights[r.Id] = r
		}
		li.IsAdmin = 1
	} else {
		var info string
		e := dbMgr.db.QueryRow("select rights from admin_account where account = ?", account).Scan(&info)
		if e != nil {
			return nil
		}
		//defer rows.Close()
		//if !rows.Next() {
		//	return nil
		//}
		//info := ""
		//e = rows.Scan(&info)
		//if e != nil {
		//	return nil
		//}

		var Extend GMInfoEx
		if len(info) != 0 {
			e = json.Unmarshal([]byte(info), &Extend)
			if e != nil {
				return nil
			}
		}
		for _, rid := range Extend.Rights {
			if rid > 0 {
				if rid < GROUP_TAG { //组
					for _, r := range AllRights {
						if r.Id/GROUP_TAG == rid {
							li.rights[r.Id] = r
						}
					}
				} else {
					for _, r := range AllRights {
						if r.Id == rid {
							li.rights[r.Id] = r
							break
						}
					}
				}
			}
		}
		for _, rid := range Extend.Rights {
			if rid < 0 {
				rid = -rid
				if _, ok := li.rights[rid]; ok {
					delete(li.rights, rid)
				}
			}
		}
	}

	return li
}

//func GetLoginInfo(req *http.Request) *LoginInfo {
//	cookie, err := req.Cookie("Name")
//	if err != nil {
//		return nil
//	}
//	if cookie.Value == "" {
//		return nil
//	}
//	token, err := req.Cookie("Token")
//	if token == nil {
//		return nil
//	}
//
//	li := &LoginInfo{Name: cookie.Value, Token: token.Value}
//	li.rights = make(map[int]*Rights)
//	if cookie.Value == DEFAULT_ADMIN {
//		for _, r := range AllRights {
//			li.rights[r.Id] = r
//		}
//		li.IsAdmin = 1
//	} else {
//		rows, e := dbMgr.db.Query("select info from admin_account where acc = ? limit 1;", cookie.Value)
//		if e != nil {
//			return nil
//		}
//		defer rows.Close()
//		if !rows.Next() {
//			return nil
//		}
//		info := ""
//		e = rows.Scan(&info)
//		if e != nil {
//			return nil
//		}
//
//		var Extend GMInfoEx
//		if len(info) != 0 {
//			e = json.Unmarshal([]byte(info), &Extend)
//			if e != nil {
//				return nil
//			}
//		}
//		for _, rid := range Extend.Rights {
//			if rid > 0 {
//				if rid < GROUP_TAG { //组
//					for _, r := range AllRights {
//						if r.Id/GROUP_TAG == rid {
//							li.rights[r.Id] = r
//						}
//					}
//				} else {
//					for _, r := range AllRights {
//						if r.Id == rid {
//							li.rights[r.Id] = r
//							break
//						}
//					}
//				}
//			}
//		}
//		for _, rid := range Extend.Rights {
//			if rid < 0 {
//				rid = -rid
//				if _, ok := li.rights[rid]; ok {
//					delete(li.rights, rid)
//				}
//			}
//		}
//	}
//
//	return li
//}

func GetGMInfoAll() (l []*GMInfo) {
	rows, e := dbMgr.db.Query("select account, password, salt, rights from admin_account;")
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		account := ""
		password := ""
		salt := ""
		rights := ""
		e = rows.Scan(&account, &password, &salt, &rights)
		if e != nil {
			return
		}

		tmp := &GMInfo{Name: account}
		if account == DEFAULT_ADMIN {
			tmp.Rights = "所有权限"
		} else if strings.Index(account, "*") == 0 {
			tmp.Rights = "账号已停用"
		} else {
			var ex GMInfoEx
			if len(rights) != 0 {
				json.Unmarshal([]byte(rights), &ex)
			}
			var rightsStr []string
			for _, v := range ex.Rights {
				rightsStr = append(rightsStr, strconv.Itoa(v))
			}
			if len(rightsStr) != 0 {
				tmp.Rights = strings.Join(rightsStr, ";")
			}
		}
		l = append(l, tmp)
	}
	return
}
func admin_set(w http.ResponseWriter, r *http.Request, l *LoginInfo) {
	//ci := GetLoginInfo(r)
	//if ci == nil || !ci.CheckRights("/admin_set") {
	//	HttpWrite(w, "web/template/404.html", nil)
	//	return
	//}
	d := make(map[string]interface{})
	d["Users"] = GetGMInfoAll()
	d["IsAdmin"] = l.IsAdmin
	d["ServerGroup"] = []int{}
	d["RightsSet"] = []int{}
	d["IsSuper"] = false
	d["RightsSetAdmin"] = "管理员"
	d["AllRights"] = AllRights
	HttpWrite(w, "web/template/adminset.html", d)
}

func CheckExistGMByName(acc string) (bool, error) {
	rows, e := dbMgr.db.Query("select count(*) from admin_account where account = ?;", acc)
	if e != nil {
		return false, e
	}
	defer rows.Close()

	if !rows.Next() {
		return false, errors.New("query rows error")
	}
	n := int(0)
	e = rows.Scan(&n)
	if e != nil {
		return false, e
	}
	return n == 1, nil
}

func AddGM(w http.ResponseWriter, r *http.Request) {
	//ci := GetLoginInfo(r)
	//if ci == nil || !ci.CheckRights("/AddGM") {
	//	HttpWriteJson(w, &CommonResult{1, "没有权限", nil})
	//	return
	//}

	r.ParseForm()
	acc := r.FormValue("Name")
	pwd := r.FormValue("Pwd")

	//logrus.Debug("AddGM:", acc, " ", pwd)
	lName := len(acc)
	if lName < 4 || lName > 20 {
		HttpWriteJson(w, &CommonResult{1, "Acc长度错误", nil})
		return
	}
	if strings.Index(acc, "*") != -1 {
		HttpWriteJson(w, &CommonResult{1, "Acc非法字符", nil})
		return
	}
	lPwd := len(pwd)
	if lPwd < 4 || lPwd > 20 {
		HttpWriteJson(w, &CommonResult{1, "Pwd长度错误", nil})
		return
	}

	exist, e := CheckExistGMByName("*" + acc)
	if e != nil {
		HttpWriteJson(w, &CommonResult{1, "添加失败：" + e.Error(), nil})
		return
	}
	if exist {
		HttpWriteJson(w, &CommonResult{1, "添加失败：已存在被停用的同名账号", nil})
		return
	}

	var salt, authMD5 string
	salt = gameutil.GetRandomString(8)
	authMD5 = gameutil.MD5(pwd + salt)
	_, e = dbMgr.db.Exec("INSERT admin_account(account, password, salt) VALUES (?, ?, ?)", acc, authMD5, salt)
	if e != nil {
		HttpWriteJson(w, &CommonResult{1, "添加失败：" + e.Error(), nil})
		return
	}

	HttpWriteJson(w, &CommonResult{0, "ok", nil})
	return
}

func UpdateGM(w http.ResponseWriter, r *http.Request) {
	//ci := GetLoginInfo(r)
	//if ci == nil || !ci.CheckRights("/UpdateGM") {
	//	HttpWriteJson(w, &CommonResult{1, "没有权限", nil})
	//	return
	//}

	r.ParseForm()
	acc := r.FormValue("Name")
	if acc == "" {
		HttpWriteJson(w, &CommonResult{0, "修改失败: 未指定账号", nil})
		return
	}
	rights := r.FormValue("Rights")

	rows, e := dbMgr.db.Query("select password, salt, rights from admin_account where account = ? limit 1;", acc)
	if e != nil {
		HttpWriteJson(w, &CommonResult{0, "修改失败:" + e.Error(), nil})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		HttpWriteJson(w, &CommonResult{0, "修改失败:没有找到账号" + e.Error(), nil})
		return
	}
	password := ""
	salt := ""
	right := ""
	e = rows.Scan(&password, &salt, &right)
	if e != nil {
		HttpWriteJson(w, &CommonResult{0, "修改失败:" + e.Error(), nil})
		return
	}

	pwd := r.FormValue("Pwd")
	if pwd != "******" {
		salt = gameutil.GetRandomString(8)
		password = gameutil.MD5(pwd + salt)
	}
	if rights == "所有权限" || rights == "账号已停用" {
		_, e = dbMgr.db.Exec("update admin_account set password=?, salt=? where account = ?;", password, salt, acc)
		if e != nil {
			HttpWriteJson(w, &CommonResult{1, "修改失败:" + e.Error(), nil})
			return
		}
	} else {
		rList, _ := gameutil.ParseInt32s(rights, ";")
		//logrus.Debug("UpdateGM:", acc, " ", rights)
		var ex GMInfoEx
		for _, r := range rList {
			ex.Rights = append(ex.Rights, int(r))
		}

		bs, e := json.Marshal(ex)
		if e != nil {
			HttpWriteJson(w, &CommonResult{0, "修改失败:" + e.Error(), nil})
			return
		}

		_, e = dbMgr.db.Exec("update admin_account set password=?, salt=?, rights = ? where account = ?;", password, salt, string(bs), acc)
		if e != nil {
			HttpWriteJson(w, &CommonResult{1, "修改失败:" + e.Error(), nil})
			return
		}
	}

	HttpWriteJson(w, &CommonResult{0, "ok", nil})
	return
}

func RemoveGM(w http.ResponseWriter, r *http.Request) {
	//ci := GetLoginInfo(r)
	//if ci == nil || !ci.CheckRights("/RemoveGM") {
	//	HttpWriteJson(w, &CommonResult{1, "没有权限", nil})
	//	return
	//}

	r.ParseForm()
	acc := r.FormValue("Name")
	if acc == DEFAULT_ADMIN {
		HttpWriteJson(w, &CommonResult{0, "停用失败: 不可停用默认账号", nil})
		return
	}
	if strings.Index(acc, "*") == 0 {
		HttpWriteJson(w, &CommonResult{0, "停用失败: 账号已停用", nil})
		return
	}
	rows, e := dbMgr.db.Query("select account from admin_account where account = ? limit 1;", acc)
	if e != nil {
		HttpWriteJson(w, &CommonResult{0, "停用失败:" + e.Error(), nil})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		HttpWriteJson(w, &CommonResult{0, "停用失败:没有找到账号" + e.Error(), nil})
		return
	}
	tmp := ""
	e = rows.Scan(&tmp)
	if e != nil {
		HttpWriteJson(w, &CommonResult{0, "停用失败:" + e.Error(), nil})
		return
	}

	_, e = dbMgr.db.Exec("update admin_account set account = ? where account = ?;", "*"+acc, acc)
	if e != nil {
		HttpWriteJson(w, &CommonResult{1, "停用失败:" + e.Error(), nil})
		return
	}

	HttpWriteJson(w, &CommonResult{0, "ok", nil})
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	//defer func() {
	//	if r := recover(); r != nil {
	//	}
	//}()

	name := r.PostFormValue("name")
	password := r.PostFormValue("password")
	if strings.Index(name, "*") != -1 {
		HttpWriteJson(w, &CommonResult{1, "登录失败:账号已停用", nil})
		return
	}
	rows, e := dbMgr.db.Query("select account, password, salt from admin_account where account = ? limit 1;", name)
	if e != nil {
		HttpWriteJson(w, &CommonResult{1, "登录失败:" + e.Error(), nil})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		HttpWriteJson(w, &CommonResult{1, "登录失败:没有找到账号", nil})
		return
	}
	account := ""
	p := ""
	salt := ""
	e = rows.Scan(&account, &p, &salt)
	if e != nil {
		HttpWriteJson(w, &CommonResult{1, "登录失败:" + e.Error(), nil})
		return
	}
	if password != "" {
		tmp := gameutil.MD5(password + salt)
		if p != tmp {
			HttpWriteJson(w, &CommonResult{1, "密码不正确", nil})
			return
		}
	}

	sess := sessionMgr.SessionStart(w, r)
	sess.Set("account", account)

	mp := make(map[string](time.Time))
	mp[name] = time.Now()
	jsonBytes, _ := json.Marshal(mp)
	ioutil.WriteFile(SessionRoot, jsonBytes, 0664)

	////cookie := http.Cookie{Name: "Token", Value: "Token", Path: "/", MaxAge: 86400}
	//http.SetCookie(res, &http.Cookie{Name: "Name", Value: name, Path: "/", MaxAge: 86400})
	//http.SetCookie(res, &http.Cookie{Name: "Token", Value: "Token", Path: "/", MaxAge: 86400})
	HttpWriteJson(w, &CommonResult{0, "", nil})
}

func Register(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	r.ParseForm()
	acc := r.PostFormValue("name")
	pwd := r.PostFormValue("password")

	//logrus.Debug("AddGM:", acc, " ", pwd)
	lName := len(acc)
	if lName < 4 || lName > 20 {
		HttpWriteJson(w, &CommonResult{1, "Acc长度错误", nil})
		return
	}
	if strings.Index(acc, "*") != -1 {
		HttpWriteJson(w, &CommonResult{1, "Acc非法字符", nil})
		return
	}
	lPwd := len(pwd)
	if lPwd < 4 || lPwd > 20 {
		HttpWriteJson(w, &CommonResult{1, "Pwd长度错误", nil})
		return
	}

	exist, e := CheckExistGMByName("*" + acc)
	if e != nil {
		HttpWriteJson(w, &CommonResult{1, "注册失败：" + e.Error(), nil})
		return
	}
	if exist {
		HttpWriteJson(w, &CommonResult{1, "注册失败：已存在被停用的同名账号", nil})
		return
	}

	var salt, authMD5 string
	salt = gameutil.GetRandomString(8)
	authMD5 = gameutil.MD5(pwd + salt)
	_, e = dbMgr.db.Exec("INSERT admin_account(account, password, salt) VALUES (?, ?, ?)", acc, authMD5, salt)
	if e != nil {
		HttpWriteJson(w, &CommonResult{1, "注册失败：" + e.Error(), nil})
		return
	}

	//sess := sessionMgr.SessionStart(w, r)
	//sess.Set("account", acc)

	//cookie := http.Cookie{Name: "Token", Value: "Token", Path: "/", MaxAge: 86400}
	//http.SetCookie(w, &http.Cookie{Name: "Name", Value: acc, Path: "/", MaxAge: 86400})
	//http.SetCookie(w, &http.Cookie{Name: "Token", Value: "Token", Path: "/", MaxAge: 86400})
	HttpWriteJson(w, &CommonResult{0, "", nil})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	//sessionMgr.SessionDestroy(w, r)
	mp := make(map[string](time.Time))
	jsonBytes, _ := json.Marshal(mp)
	ioutil.WriteFile(SessionRoot, jsonBytes, 0664)

	//http.SetCookie(res, &http.Cookie{Name: "Name", Value: "", Path: "/", MaxAge: 1})
	//http.SetCookie(res, &http.Cookie{Name: "Token", Value: "", Path: "/", MaxAge: 1})
	HttpWriteJson(w, &CommonResult{1, "", nil})
}

func Password(w http.ResponseWriter, r *http.Request) {
	//ci := GetLoginInfo(r)
	//if ci == nil {
	//	HttpWrite(w, "web/template/404.html", nil)
	//	return
	//}
	d := make(map[string]interface{})
	HttpWrite(w, "web/template/password.html", d)
}

type PasswordEditForm struct {
	Name         string
	PasswordOld  string
	PasswordNew1 string
	PasswordNew2 string
}

func PasswordEdit(w http.ResponseWriter, r *http.Request, account string) {
	//ci := GetLoginInfo(r)
	//if ci == nil {
	//	HttpWriteJson(w, &CommonResult{1, "没有登录", nil})
	//	return
	//}

	r.ParseForm()
	reqInfo := PasswordEditForm{}

	reqInfo.Name = account //r.PostFormValue("Name")
	reqInfo.PasswordOld = r.PostFormValue("PasswordOld")
	reqInfo.PasswordNew1 = r.PostFormValue("PasswordNew1")
	reqInfo.PasswordNew2 = r.PostFormValue("PasswordNew2")

	if reqInfo.Name == "" || reqInfo.PasswordNew1 == "" || reqInfo.PasswordNew1 != reqInfo.PasswordNew2 {
		HttpWriteJson(w, &CommonResult{1, "错误信息", nil})
		return
	}

	rows, e := dbMgr.db.Query("select account, password, salt from admin_account where account = ? limit 1;", reqInfo.Name)
	if e != nil {
		HttpWriteJson(w, &CommonResult{1, e.Error(), nil})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		HttpWriteJson(w, &CommonResult{1, "没有找到信息", nil})
		return
	}
	acc := ""
	auth := ""
	salt := ""
	e = rows.Scan(&acc, &auth, &salt)
	if e != nil {
		HttpWriteJson(w, &CommonResult{1, e.Error(), nil})
		return
	}
	if auth != "" {
		if reqInfo.PasswordOld == "" {
			HttpWriteJson(w, &CommonResult{1, "密码不正确", nil})
			return
		}
		tmp := gameutil.MD5(reqInfo.PasswordOld + salt)
		if auth != tmp {
			HttpWriteJson(w, &CommonResult{1, "密码不正确", nil})
			return
		}
	}
	{
		var salt, authMD5 string
		salt = gameutil.GetRandomString(8)
		authMD5 = gameutil.MD5(reqInfo.PasswordNew1 + salt)
		_, e = dbMgr.db.Exec("update admin_account set password=?, salt=? where account=?;", authMD5, salt, acc)
		if e != nil {
			HttpWriteJson(w, &CommonResult{1, e.Error(), nil})
			return
		}
	}

	HttpWriteJson(w, &CommonResult{0, "ok", nil})
}
