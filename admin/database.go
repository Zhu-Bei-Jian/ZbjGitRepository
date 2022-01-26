package admin

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/db"
	"strings"
	"time"
)

var ErrCanNotFind = errors.New("ErrCanNotFind")

type database struct {
	db *sql.DB
}

func openDB(DBGame string) (*database, error) {
	db, err := gameutil.OpenDB(DBGame)
	if err != nil {
		return nil, errors.New("open:" + DBGame + " err:" + err.Error())
	}

	return &database{
		db: db,
	}, nil
}

func (p *database) close() {
	p.db.Close()
}

func (p *database) GetServerInfo() (map[string]string, error) {
	rows, e := dbMgr.db.Query("select `key`,`value` FROM `admin_server_info`;")
	if e != nil {
		return nil, e
	}
	if rows == nil {
		return nil, errors.New("GetServerInfo rows == nil")
	}
	defer rows.Close()

	infoMap := make(map[string]string)
	for rows.Next() {
		var key string
		var value string
		e := rows.Scan(&key, &value)
		if e == nil {
			infoMap[key] = value
		}
	}
	return infoMap, nil
}

func (p *database) UpdateServerInfo(key, value string) error {
	_, e := dbMgr.db.Exec("INSERT INTO admin_server_info(`key`,`value`)VALUES(?,?) ON DUPLICATE KEY UPDATE `value` = ?;", key, value, value)
	if e != nil {
		return e
	}
	return nil
}

func (p *database) QueryRawUserData(userID uint64, columns string, f func(rows *sql.Rows, err error)) {
	tableName := gameshared.GetUserTableName(userID)
	if columns == "" {
		columns = "*"
	}
	rows, err := p.db.Query(fmt.Sprintf("select %s from %s where userid = ?", columns, tableName), userID)
	if err != nil {
		f(nil, err)
		return
	}
	defer rows.Close()
	f(rows, err)
}

func (p *database) QueryAllRawUserDataWithMd5(userID uint64, logToFile bool) (content string, md5 string, e error) {
	p.QueryRawUserData(uint64(userID), "*", func(rows *sql.Rows, err error) {
		if err != nil {
			e = err
			return
		}
		columns, err := rows.ColumnTypes()
		if err != nil {
			e = err
			return
		}

		if !rows.Next() {
			e = errors.New("no rows")
			return

		}
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)
		if err != nil {
			e = err
			return
		}
		var columnsStringList []string
		for i, pval := range values {
			tmp := DBColumnInfo(columns[i], pval)
			columnsStringList = append(columnsStringList, tmp)
		}
		content = strings.Join(columnsStringList, "\n")
		if logToFile {
			fpath := fmt.Sprintf("./raw/%s_%d.txt", time.Now().Format("20060102150405"), userID)
			err = os.MkdirAll(path.Dir(fpath), 0766)
			if err != nil {
				logrus.Error("QueryRawUserData Save " + fpath + " " + err.Error())
			} else {
				//ioutil.WriteFile(fpath, []byte(content), 0666)
				f, _ := os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666) //打开文件
				if f != nil {
					defer f.Close()
					f.Write([]byte(content))
					//io.WriteString(f, content)
				}
			}
		}
		md5 = gameutil.MD5(content)
		return
	})
	return
}

func DBColumnInfo(column *sql.ColumnType, pval interface{}) string {
	var s_txt string
	switch v := pval.(type) {
	case nil:
		s_txt = "NULL"
	case time.Time:
		s_txt = "'" + v.Format("2006-01-02 15:04:05") + "'"
	case int, int8, int16, int32, int64, float32, float64, byte:
		s_txt = fmt.Sprint(v)
	case []byte:
		s_txt = fmt.Sprintf("%d %x", len(v), v)
	case bool:
		if v {
			s_txt = "'1'"
		} else {
			s_txt = "'0'"
		}
	default:
		s_txt = "'" + fmt.Sprint(v) + "'"
	}
	tmp := fmt.Sprint(column.Name(), " ", column.ScanType().Name(), " ", s_txt)
	return tmp
}

func ColumnProtoMessage(column string) (content proto.Message) {
	switch column {
	case "char_info":
		content = &db.CharInfo{}

	default:
	}
	return
}
