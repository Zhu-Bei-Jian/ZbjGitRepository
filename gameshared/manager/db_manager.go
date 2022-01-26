package manager

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"sanguosha.com/baselib/ioservice"
)

const (
	DB_CONNECTIONS = 10
)

type dbManager struct {
	db     *sql.DB
	worker []ioservice.IOService
}

func NewDB(source string) (*dbManager, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	d := &dbManager{
		db: db,
	}
	d.worker = make([]ioservice.IOService, DB_CONNECTIONS)
	for i := 0; i < DB_CONNECTIONS; i++ {
		d.worker[i] = ioservice.NewIOService(fmt.Sprintf("db-worker-%d-%d", 100, i), 10240)
		d.worker[i].Init()
	}
	return d, nil
}

func (p *dbManager) Run() {
	for _, v := range p.worker {
		v.Run()
	}
}

func (p *dbManager) Post(id interface{}, f func()) bool {
	i := GetHashByID(id, DB_CONNECTIONS)
	p.worker[i].Post(f)
	return true
}

func GetHashByID(id interface{}, num uint32) uint32 {
	pos := crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v", id))) % num
	return pos
}

func (p *dbManager) Close() {
	for _, v := range p.worker {
		//logrus.Infof("Database.Close... IOServiceIdx: %d", i)
		v.Fini()
		//Log("Database.Close OK, IOServiceIdx: %d", i)
	}
	p.db.Close()
}

func (p *dbManager) DB() *sql.DB {
	return p.db
}
