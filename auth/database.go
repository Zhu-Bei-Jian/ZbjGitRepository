package auth

import (
	"database/sql"
	"sanguosha.com/sgs_herox/gameutil"
	"strconv"
	"time"
)

type database struct {
	db *sql.DB
}

func NewDatabase(source string) (*database, error) {
	db, err := gameutil.OpenDB(source)
	//var key, value string
	//db.QueryRow("select * from admin_server_info ").Scan(&key, &value)
	if err != nil {
		return nil, err
	}
	d := &database{
		db: db,
	}
	return d, nil
}

func (p *database) close() {
	p.db.Close()
}

func (d *database) getAccountUserID(account string) (userid uint64, err error) {
	err = d.db.QueryRow("SELECT userid FROM account WHERE account=?", account).Scan(&userid)
	if err != nil {
		return 0, err
	}
	return userid, nil
}

func (d *database) createAccount(account string, accountType int, registerIP string, loginId string, nickname string, avatarUrl string, avatarFrameUrl string, gender int32, originData string) (uint64, error) {

	now := time.Now()

	tx, err := d.db.Begin()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var key, value string
	d.db.QueryRow("select * from admin_server_info ").Scan(&key, &value)
	//ret, err := tx.Exec("INSERT INTO user(level,login_time,logout_time)VALUES(1,?,?)", now, now)
	ret, err := d.db.Exec("INSERT INTO user(level,login_time,logout_time)VALUES(1,?,?)", now, now)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	userId, err := ret.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sName := UnicodeEmojiCode(nickname)
	_, err = tx.Exec("INSERT INTO account(userid,account_type,account,loginid,head_img_url,headframe_img_url,nickname,sex,born_time,created_ip,origin_data) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		userId, accountType, account, loginId, avatarUrl, avatarFrameUrl, sName, gender, now.Unix(), registerIP, originData)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return uint64(userId), nil
}

func (d *database) updateLoginInfo(userid uint64, loginIP string) error {
	_, err := d.db.Exec("UPDATE account SET login_ip =?,login_time=? WHERE userId=?", loginIP, time.Now(), userid)
	return err
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------
// 名字中的表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}
