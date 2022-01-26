package gameutil

import (
	"database/sql"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

func OpenDB(source string) (*sql.DB, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * time.Duration(5400))
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func OpenDBWithMaxConn(source string, maxConn int) (*sql.DB, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxConn / 2)
	//db.SetConnMaxLifetime(time.Second * time.Duration(86400))
	db.SetConnMaxLifetime(time.Second * time.Duration(5400))
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func OpenDBWithMaxMinConn(source string, maxConn, minConn int) (*sql.DB, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(minConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(5400))
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
