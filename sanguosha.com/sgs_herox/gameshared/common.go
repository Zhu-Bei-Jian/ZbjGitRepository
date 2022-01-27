package gameshared

import (
	"fmt"
	"math/rand"
	"sanguosha.com/sgs_herox/gameutil"
	"time"
)

const UserIDBase uint64 = 1000000000 //这个值表示每张分表的最大角色数 9.999E
const UserTableIDRobot uint64 = 99   //TODO user99 理解成robot的表，不真实存在 <<<假定user分表不会到99，如需超过 修改这个值>>>
const UserTableIDLobbyRobot uint64 = 98
const TotalUserTable = 10

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewShowID() int64 {
	vRand0 := rand.Int63()
	vRand1 := rand.Int63()
	vRand2 := rand.Int63()

	n1 := vRand0%9 + 1
	//n23 := vRand0 % 100
	n47 := vRand1 % 10000
	n8e := vRand2 % 10000
	//n47 := int64(1314) + vRand1%8600
	//n8e := int64(173) + vRand2%9600

	//showid := n1*10000000000 + n23*100000000 + n47*10000 + n8e
	showid := n1*100000000 + n47*10000 + n8e

	return showid
}

func NewUserID(tableId int64) int64 {
	vRand0 := rand.Int63()
	vRand1 := rand.Int63()
	vRand2 := rand.Int63()
	n1 := vRand0%9 + 1
	n47 := vRand1 % 10000
	n8e := vRand2 % 1000
	sign := (tableId*n1*24357 + n47 + n8e) % 10
	userid := tableId*1000000000 + n1*100000000 + n47*10000 + n8e*10 + sign
	return userid
}

func IsUserIDOK(uid uint64) bool {
	userTableID := GetUserTableID(uid)
	if userTableID == 0 {
		return false
	}
	return userTableID <= TotalUserTable
}

func RandUserTable() (string, int64) {
	tableID := gameutil.Rand63()%TotalUserTable + 1
	return fmt.Sprintf("user%d", tableID), tableID
}

func GetUserTableID(uid uint64) int64 {
	return int64(uid / UserIDBase)
}

func GetUserTableName(uid uint64) string {
	return fmt.Sprintf("user%d", GetUserTableID(uid))
}

func GetGoldTableName(uid uint64) string {
	return fmt.Sprintf("gold%d", GetUserTableID(uid))
}
