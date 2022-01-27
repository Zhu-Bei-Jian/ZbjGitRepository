package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

var srcAccount = flag.String("srcAccount", "", "source account")
var dstAccount = flag.String("dstAccount", "", "destination account")

func main() {
	flag.Parse()

	fmt.Printf("src_account:%s dst_account:%s \n", *srcAccount, *dstAccount)
	fmt.Println("start copy...")

	data, err := ioutil.ReadFile("config.json")
	checkErr(err)
	var cfg map[string]string
	err = json.Unmarshal(data, &cfg)
	checkErr(err)
	dbSourceSrc := cfg["src_ds"]
	dbSourceDst := cfg["dst_ds"]

	srcDB, err := sql.Open("mysql", dbSourceSrc)
	checkErr(err)
	defer srcDB.Close()

	dstDB, err := sql.Open("mysql", dbSourceDst)
	checkErr(err)
	defer dstDB.Close()

	var userId int64
	err = srcDB.QueryRow("select userid from auth where account = ?", *srcAccount).Scan(&userId)
	if err == sql.ErrNoRows {
		log.Fatalf(fmt.Sprintf("srcAccount(%s) not exist", *srcAccount))
	}
	checkErr(err)

	uid, err := insert(userId, srcDB, dstDB)
	checkErr(err)

	_, err = dstDB.Exec("update auth set userid=? where account = ? ", uid, *dstAccount)
	checkErr(err)

	fmt.Println("copy success!")
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func NewUserID(tableId int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	vRand0 := r.Int63()
	vRand1 := r.Int63()
	vRand2 := r.Int63()
	n1 := vRand0%9 + 1
	n47 := vRand1 % 10000
	n8e := vRand2 % 1000
	sign := (tableId*n1*24357 + n47 + n8e) % 10
	userid := tableId*1000000000 + n1*100000000 + n47*10000 + n8e*10 + sign
	return userid
}

func insert(uid int64, srcDB *sql.DB, dstDB *sql.DB) (int64, error) {
	tableId := GetTableId(uid)
	rows, errDB1Q := srcDB.Query(fmt.Sprintf("select * from user%d where userid=%d", tableId, uid))
	defer rows.Close()
	if errDB1Q != nil {
		return 0, errDB1Q
	}
	cols, errCol := rows.Columns()
	if errCol != nil {
		return 0, errCol
	}
	rawRes := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols))
	var fields string
	for i := 0; i < len(cols); i++ {
		dest[i] = &rawRes[i]
		if i > 0 {
			fields += cols[i] + ","
		}
	}
	fields = fields[:len(fields)-1]
	placeHolder := "?"
	for i := 0; i < len(cols)-2; i++ {
		placeHolder = placeHolder + ",?"
	}
	for rows.Next() {
		rows.Scan(dest...)
		break
	}
	query := fmt.Sprintf("insert into user%d(%s) values(%s)", tableId, fields, placeHolder)
	res, err1 := dstDB.Exec(query, dest[1:]...)
	if err1 != nil {
		return 0, err1
	}
	nowUid, err2 := res.LastInsertId() //自增生成的 userID ,
	newUid := nowUid
	if err2 != nil {
		return 0, err2
	}
	var tempUid int64
	for i := 0; i < 100; i++ {
		tempUid = NewUserID(tableId)
		if !IsContainsKey(tempUid, dstDB) {
			newUid = tempUid
			break
		}
	}
	_, err3 := dstDB.Exec(fmt.Sprintf("update user%v set userid=%v where userid=%v", tableId, newUid, nowUid))
	if err3 != nil {
		return 0, err3
	}
	return newUid, nil
}

func GetTableId(uid int64) int64 {
	return uid / 1000000000
}

func IsContainsKey(uid int64, db *sql.DB) bool { // 数据库中查找 是否 userid 已存在
	var cnt int

	db.QueryRow(fmt.Sprintf("select count(userid) from user%v where userid =%v", GetTableId(uid), uid)).Scan(&cnt)
	return cnt == 1
}
