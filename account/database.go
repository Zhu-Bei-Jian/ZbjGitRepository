package account

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

var ErrCanNotFind = errors.New("ErrCanNotFind")

type database struct {
	db *sql.DB
}

func newDatabase(source string) (*database, error) {
	db, err := gameutil.OpenDB(source)
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

func (p *database) FindUserIdByShowID(showID uint64) (uint64, error) {
	var userId uint64
	err := p.db.QueryRow("SELECT userid FROM account where showID = ?", showID).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrCanNotFind
		}
		return 0, err
	}
	return userId, nil
}

func (p *database) FindUserIdByAccount(account string) (uint64, error) {
	var userId uint64
	err := p.db.QueryRow("SELECT userid FROM account where account = ?", account).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrCanNotFind
		}
		return 0, err
	}
	return userId, nil
}

func (p *database) FindUserIdByUnionId(unionId string, accountType gameconf.AccountLoginTyp) (uint64, error) {
	var userId uint64
	err := p.db.QueryRow("SELECT userid FROM account where unionid = ? and account_type=?", unionId, int32(accountType)).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrCanNotFind
		}
		return 0, err
	}
	return userId, nil
}

func (d *database) FindUserIdsByNickName(nickName string) ([]uint64, error) {
	rows, err := d.db.Query("SELECT userid FROM account WHERE nickname=?", nickName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIds []uint64
	for rows.Next() {
		var userID uint64
		err := rows.Scan(&userID)
		if err == nil {
			userIds = append(userIds, userID)
		}
	}
	return userIds, nil
}

func (d *database) FindUserIdsLikeNickName(nickName string) ([]uint64, error) {
	rows, err := d.db.Query("SELECT userid FROM account WHERE nickname like ?", "%"+nickName+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIds []uint64
	for rows.Next() {
		var userID uint64
		err := rows.Scan(&userID)
		if err == nil {
			userIds = append(userIds, userID)
		}
	}
	return userIds, nil
}
